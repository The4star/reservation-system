package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/the4star/reservation-system/pkg/config"
	"github.com/the4star/reservation-system/pkg/models"
	"github.com/the4star/reservation-system/pkg/render"
)

// the repository used by the handlers
var Repo *Repository

// the repository type
type Repository struct {
	App *config.AppConfig
}

//creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// Set the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home renders the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About renders the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{})
}

// StandardSuite renders the standard suite page
func (m *Repository) StandardSuite(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "standard-suite.page.tmpl", &models.TemplateData{})
}

// DeluxeSuite renders the deluxe suite page
func (m *Repository) DeluxeSuite(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "deluxe-suite.page.tmpl", &models.TemplateData{})
}

// Availability renders the availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "availability.page.tmpl", &models.TemplateData{})
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
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Book renders the book page
func (m *Repository) Book(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "book.page.tmpl", &models.TemplateData{})
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}
