package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/the4star/reservation-system/internal/config"
	"github.com/the4star/reservation-system/internal/driver"
	"github.com/the4star/reservation-system/internal/forms"
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

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo creates a new repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// NewHandlers Set the repository for the handlers
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

// ChooseRoom sets the room the user has chosen and redirects to the book page
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	splitURI := strings.Split(r.URL.Path, "/")
	if len(splitURI) != 3 {
		m.App.Session.Put(r.Context(), "error", "Error accessing room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	roomID, err := strconv.Atoi(splitURI[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error accessing room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Error accessing room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/book", http.StatusSeeOther)
}

// PostAvailability handles the availability page form
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error retrieving form information")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	start := r.Form.Get("start-date")
	end := r.Form.Get("end-date")
	timeLayout := "2006-01-02"
	startDate, err := time.Parse(timeLayout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error processing dates")
		http.Redirect(w, r, "/availability", http.StatusSeeOther)
		return
	}
	endDate, err := time.Parse(timeLayout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error processing dates")
		http.Redirect(w, r, "/availability", http.StatusSeeOther)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error processing request")
		http.Redirect(w, r, "/availability", http.StatusSeeOther)
		return
	}

	if len(rooms) == 0 {
		m.App.InfoLog.Println("No availablity")
		m.App.Session.Put(r.Context(), "error", "No availability for selected dates")
		http.Redirect(w, r, "/availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

type RoomAvailabilityRequest struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	RoomID    string `json:"roomID"`
}
type roomAvailabilityResponse struct {
	OK  bool   `json:"ok"`
	Msg string `json:"msg"`
}

// RoomAvailability handles the form on room pages
func (m *Repository) PostRoomAvailability(w http.ResponseWriter, r *http.Request) {
	jsonData := RoomAvailabilityRequest{}
	data, _ := io.ReadAll(r.Body)

	json.Unmarshal(data, &jsonData)

	fmt.Printf("%+v", jsonData)
	w.Header().Set("Content-Type", "application/json")

	timeLayout := "2006-01-02"
	startDate, err := time.Parse(timeLayout, jsonData.StartDate)
	if err != nil {
		internalServerErrorJSON(w)
		return
	}
	endDate, err := time.Parse(timeLayout, jsonData.EndDate)
	if err != nil {
		internalServerErrorJSON(w)
		return
	}

	roomID, err := strconv.Atoi(jsonData.RoomID)
	if err != nil {
		internalServerErrorJSON(w)
		return
	}

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(roomID, startDate, endDate)
	if err != nil {
		internalServerErrorJSON(w)
		return
	}

	resp := roomAvailabilityResponse{
		OK: available,
	}

	responseData, err := json.Marshal(resp)
	if err != nil {
		internalServerErrorJSON(w)
		return
	}

	w.Write(responseData)
}

// BookRoom renders the booking page using params
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error processing request")
		http.Redirect(w, r, "/availability", http.StatusTemporaryRedirect)
		return
	}
	sd := r.URL.Query().Get("sd")
	ed := r.URL.Query().Get("ed")

	timeLayout := "2006-01-02"
	startDate, err := time.Parse(timeLayout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error processing request")
		http.Redirect(w, r, "/availability", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(timeLayout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error processing request")
		http.Redirect(w, r, "/availability", http.StatusTemporaryRedirect)
		return
	}
	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error processing request")
		http.Redirect(w, r, "/availability", http.StatusTemporaryRedirect)
		return
	}

	var res models.Reservation

	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate
	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/book", http.StatusSeeOther)
}

// Book renders the book page
func (m *Repository) Book(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Error retrieving reservation")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error retrieving room id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	res.Room.RoomName = room.RoomName
	m.App.Session.Put(r.Context(), "reservation", res)

	stringMap := make(map[string]string)
	stringMap["start-date"] = res.StartDate.Format("2006-01-02")
	stringMap["end-date"] = res.EndDate.Format("2006-01-02")

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "book.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Form:      forms.New(nil),
		Data:      data,
	})
}

// PostBook handles the booking form
func (m *Repository) PostBook(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Error retrieving reservation")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error retrieving form information")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("first-name", "last-name", "email", "room-id")
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

	reservation.FirstName = r.Form.Get("first-name")
	reservation.LastName = r.Form.Get("last-name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

	//save to db
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error creating reservation")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.ID = newReservationID
	m.App.Session.Put(r.Context(), "reservation", reservation)

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Error creating reservation restriction")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
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

	stringMap := make(map[string]string)
	stringMap["start-date"] = reservation.StartDate.Format("2006-01-02")
	stringMap["end-date"] = reservation.EndDate.Format("2006-01-02")

	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
	})
}
