package handlers

import (
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

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{})
}

func (m *Repository) StandardSuite(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "standard-suite.page.tmpl", &models.TemplateData{})
}

func (m *Repository) DeluxeSuite(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "deluxe-suite.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "availability.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Book(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "book.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contact.page.tmpl", &models.TemplateData{})
}
