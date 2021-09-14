package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/the4star/reservation-system/internal/config"
	"github.com/the4star/reservation-system/internal/driver"
	"github.com/the4star/reservation-system/internal/forms"
	"github.com/the4star/reservation-system/internal/helpers"
	"github.com/the4star/reservation-system/internal/models"
	"github.com/the4star/reservation-system/internal/render"
	"github.com/the4star/reservation-system/internal/repository"
	"github.com/the4star/reservation-system/internal/repository/dbrepo"
)

// the repository used by the handlers
var Repo *Repository

// the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

//creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// Set the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home renders the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About renders the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// StandardSuite renders the standard suite page
func (m *Repository) StandardSuite(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "standard-suite.page.tmpl", &models.TemplateData{})
}

// DeluxeSuite renders the deluxe suite page
func (m *Repository) DeluxeSuite(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "deluxe-suite.page.tmpl", &models.TemplateData{})
}

// Availability renders the availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability handles the availability page form
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start-date")
	end := r.Form.Get("end-date")

	fmt.Println(start, end)
	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type roomAvailabilityRequest struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	RoomType  string `json:"roomType"`
}
type roomAvailabilityResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// RoomAvailability handles the form on room pages
func (m *Repository) RoomAvailability(w http.ResponseWriter, r *http.Request) {
	jsonData := roomAvailabilityRequest{}
	data, _ := io.ReadAll(r.Body)

	json.Unmarshal(data, &jsonData)

	fmt.Printf("%+v", jsonData)

	resp := roomAvailabilityResponse{
		OK:      true,
		Message: "Available",
	}

	out, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Book renders the book page
func (m *Repository) Book(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.Template(w, r, "book.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostBook handles the booking form
func (m *Repository) PostBook(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd := r.Form.Get("start-date")
	ed := r.Form.Get("end-date")

	timeLayout := "2006-01-02"
	startDate, err := time.Parse(timeLayout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate, err := time.Parse(timeLayout, ed)
	if err != nil {
		helpers.ServerError(w, err)
	}

	RoomID, err := strconv.Atoi(r.Form.Get("room-id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    RoomID,
	}

	form := forms.New(r.PostForm)
	form.Required("first-name", "last-name", "email")
	form.MinLength("first-name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		fmt.Println(form)
		render.Template(w, r, "book.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	//save to db
	err = m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
	}
	//redirect to summary if form validation passes
	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get reservation from session.")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
