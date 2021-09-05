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
	"github.com/the4star/reservation-system/internal/models"
	"github.com/the4star/reservation-system/internal/render"
)

var app config.AppConfig
var functions = template.FuncMap{}
var session *scs.SessionManager
var pathToTemplates string = "../../templates"

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func getRoutes() http.Handler {
	gob.Register(models.Reservation{})

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	templateCache, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache.")
	}

	app.TemplateCache = templateCache
	app.UseCache = true

	repo := NewRepo(&app)
	NewHandlers(repo)
	render.NewTemplates(&app)

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(SessionLoad)

	noSurfGroup := router.Group(nil)
	// noSurfGroup.Use(NoSurf)

	//routes
	router.Get("/", Repo.Home)
	router.Get("/about", Repo.About)
	router.Get("/contact", Repo.Contact)
	router.Get("/rooms/standard-suite", Repo.StandardSuite)
	router.Get("/rooms/deluxe-suite", Repo.DeluxeSuite)

	router.Get("/book", Repo.Book)
	router.Post("/book", Repo.PostBook)
	router.Get("/reservation-summary", Repo.ReservationSummary)

	noSurfGroup.Get("/availability", Repo.Availability)
	noSurfGroup.Post("/availability", Repo.PostAvailability)
	router.Post("/room-availability", Repo.RoomAvailability)

	fileServer := http.FileServer(http.Dir("./static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return router
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
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}

	return myCache, nil
}
