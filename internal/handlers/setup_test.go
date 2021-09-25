package handlers

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"text/template"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
	"github.com/the4star/reservation-system/internal/config"
	"github.com/the4star/reservation-system/internal/helpers"
	"github.com/the4star/reservation-system/internal/models"
	"github.com/the4star/reservation-system/internal/render"
)

var app config.AppConfig
var functions = template.FuncMap{
	"niceDate":   helpers.NiceDate,
	"formatDate": helpers.FormatDate,
	"iterate":    helpers.Iterate,
}
var session *scs.SessionManager
var pathToTemplates string = "../../templates"
var infoLog *log.Logger
var errorLog *log.Logger

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	app.InProduction = false

	// create loggers
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	defer close(mailChan)
	listenForMail()

	templateCache, err := CreateTestTemplateCache()
	if err != nil {
		log.Println(err)
		log.Fatal("cannot create template cache.")
	}

	app.TemplateCache = templateCache
	app.UseCache = true

	repo := NewTestRepo(&app)
	NewHandlers(repo)
	render.NewRenderer(&app)
	os.Exit(m.Run())
}

func getRoutes() http.Handler {
	helpers.NewHelpers(&app)
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(SessionLoad)

	//routes
	router.Get("/", Repo.Home)
	router.Get("/about", Repo.About)
	router.Get("/contact", Repo.Contact)
	router.Get("/rooms/standard-suite", Repo.StandardSuite)
	router.Get("/rooms/deluxe-suite", Repo.DeluxeSuite)

	//booking
	router.Get("/book", Repo.Book)
	router.Post("/book", Repo.PostBook)
	router.Get("/reservation-summary", Repo.ReservationSummary)

	//user
	router.Get("/user/login", Repo.ShowLogin)
	router.Post("/user/login", Repo.PostLogin)
	router.Get("/user/logout", Repo.Logout)

	//protected
	router.Get("/admin/dashboard", Repo.AdminDashboard)
	router.Get("/admin/reservations-new", Repo.AdminNewReservations)
	router.Get("/admin/reservations-all", Repo.AdminAllReservations)
	router.Get("/admin/reservations-calendar", Repo.AdminReservationsCalendar)
	router.Post("/admin/reservations-calendar", Repo.AdminPostReservationsCalendar)
	router.Get("/admin/reservations/{src}/{id}/show", Repo.AdminShowReservation)
	router.Post("/admin/reservations/{src}/{id}", Repo.AdminPostUpdateReservation)
	router.Get("/admin/reservations/{src}/{id}/process", Repo.AdminProcessReservation)
	router.Get("/admin/reservations/{src}/{id}/delete", Repo.AdminDeleteReservation)

	router.Get("/availability", Repo.Availability)
	router.Post("/availability", Repo.PostAvailability)
	router.Post("/room-availability", Repo.PostRoomAvailability)

	fileServer := http.FileServer(http.Dir("./static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return router
}

func listenForMail() {
	go func() {
		for {
			<-app.MailChan
		}
	}()
}

// adds csrf protection to all POST requests.
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: false,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// creates a template cache as a map.
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	childPages, err := filepath.Glob(fmt.Sprintf("%s/*/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	pages = append(pages, childPages...)

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/layouts/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/layouts/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}

	return myCache, nil
}
