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
		m.App.ErrorLog.Println(err)
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
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error retrieving form information")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	start := r.Form.Get("start-date")
	end := r.Form.Get("end-date")
	timeLayout := "2006-01-2"
	startDate, err := time.Parse(timeLayout, start)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error processing dates")
		http.Redirect(w, r, "/availability", http.StatusSeeOther)
		return
	}
	endDate, err := time.Parse(timeLayout, end)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error processing dates")
		http.Redirect(w, r, "/availability", http.StatusSeeOther)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.ErrorLog.Println(err)
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

	timeLayout := "2006-01-2"
	startDate, err := time.Parse(timeLayout, jsonData.StartDate)
	if err != nil {
		m.App.ErrorLog.Println(err)
		internalServerErrorJSON(w)
		return
	}
	endDate, err := time.Parse(timeLayout, jsonData.EndDate)
	if err != nil {
		m.App.ErrorLog.Println(err)
		internalServerErrorJSON(w)
		return
	}

	roomID, err := strconv.Atoi(jsonData.RoomID)
	if err != nil {
		m.App.ErrorLog.Println(err)
		internalServerErrorJSON(w)
		return
	}

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(roomID, startDate, endDate)
	if err != nil {
		m.App.ErrorLog.Println(err)
		internalServerErrorJSON(w)
		return
	}

	resp := roomAvailabilityResponse{
		OK: available,
	}

	responseData, err := json.Marshal(resp)
	if err != nil {
		m.App.ErrorLog.Println(err)
		internalServerErrorJSON(w)
		return
	}

	w.Write(responseData)
}

// BookRoom renders the booking page using params
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error processing request")
		http.Redirect(w, r, "/availability", http.StatusTemporaryRedirect)
		return
	}
	sd := r.URL.Query().Get("sd")
	ed := r.URL.Query().Get("ed")

	timeLayout := "2006-01-2"
	startDate, err := time.Parse(timeLayout, sd)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error processing request")
		http.Redirect(w, r, "/availability", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(timeLayout, ed)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error processing request")
		http.Redirect(w, r, "/availability", http.StatusTemporaryRedirect)
		return
	}
	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.ErrorLog.Println(err)
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
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error retrieving room id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	res.Room.RoomName = room.RoomName
	m.App.Session.Put(r.Context(), "reservation", res)

	stringMap := make(map[string]string)
	stringMap["start-date"] = res.StartDate.Format("2006-01-2")
	stringMap["end-date"] = res.EndDate.Format("2006-01-2")

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
		m.App.ErrorLog.Println(err)
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
		m.App.ErrorLog.Println(err)
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
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error creating reservation restriction")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//send notifications
	guestMessage := fmt.Sprintf(`
		<p><strong>Reservation Confirmed</strong></p>
		<p>Dear %s,<p>
		<p>This is to confirm your reservation for the %s from %s to %s.</p>
	`,
		reservation.FirstName,
		reservation.Room.RoomName,
		helpers.NiceDate(reservation.StartDate),
		helpers.NiceDate(reservation.EndDate),
	)

	hotelMessage := fmt.Sprintf(`
		<p><strong>Reservation Confirmed</strong></p>
		<p>Dear Owner,</p>
		<p>This is to confirm a new reservation in the %s from %s to %s.</p>
	`,
		reservation.Room.RoomName,
		helpers.NiceDate(reservation.StartDate),
		helpers.NiceDate(reservation.EndDate),
	)

	sendEmail(
		m.App,
		reservation.Email,
		"mail@4starOnRegent.com",
		"Reservation Confirmation",
		guestMessage,
		"basic.html",
	)

	sendEmail(
		m.App,
		"owner@4staronregent.com",
		"mail@4staronregent.com",
		"Reservation Confirmation",
		hotelMessage,
		"basic.html",
	)

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
	stringMap["start-date"] = helpers.NiceDate(reservation.StartDate)
	stringMap["end-date"] = helpers.NiceDate(reservation.EndDate)

	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
	})
}

// ShowLogin renders the login page
func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostLogin handles logging the user in
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())
	err := r.ParseForm()
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error processing login request")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid login details")
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

// Logout logs a user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	m.App.Session.Put(r.Context(), "flash", "Successfully logged out")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// AdminDashboard shows the admin dashboard.
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

// AdminNewReservations shows the admin new reservations page.
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.GetAllNewReservations()
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error getting all new reservations")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.Template(w, r, "admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminAllReservations shows the admin all reservations page.
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.GetAllReservations()
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error getting all reservations")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.Template(w, r, "admin-all-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminShowReservation shows a single reservation
func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	splitURI := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(splitURI[4])
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error getting reservation")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	src := splitURI[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src
	stringMap["month"] = month
	stringMap["year"] = year

	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error getting reservation")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "admin-reservations-show.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
		Form:      forms.New(nil),
	})
}

// AdminPostUpdateReservation updates a single reservation
func (m *Repository) AdminPostUpdateReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error processing Form")
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		return
	}

	splitURI := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(splitURI[4])
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error getting reservation")
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		return
	}

	src := splitURI[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src

	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error getting reservation")
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("first-name", "last-name", "email", "phone")
	form.IsEmail("email")
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = res
		render.Template(w, r, "admin-reservations-show.page.tmpl", &models.TemplateData{
			Data:      data,
			StringMap: stringMap,
			Form:      form,
		})
		return
	}

	res.FirstName = r.Form.Get("first-name")
	res.LastName = r.Form.Get("last-name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	err = m.DB.UpdateReservation(res)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error updating reservation")
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		return
	}

	month := r.Form.Get("month")
	year := r.Form.Get("year")

	m.App.Session.Put(r.Context(), "flash", fmt.Sprintf("Changes saved for reservation #%d", id))

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}

}

// AdminProcessReservation updates a reservation to processed.
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	splitURI := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(splitURI[4])
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error getting reservation")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	src := splitURI[3]

	err = m.DB.UpdateProcessedForReservation(id, true)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error getting reservation")
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations/%s/%d", src, id), http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Reservation Marked as processed")

	if src == "cal" {
		year := r.URL.Query().Get("y")
		month := r.URL.Query().Get("m")
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	}
}

// AdminDeleteReservation deletes a reservation.
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	splitURI := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(splitURI[4])
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error getting reservation")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	src := splitURI[3]

	err = m.DB.DeleteReservation(id)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error getting reservation")
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations/%s/%d", src, id), http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "Reservation Deleted")

	if src == "cal" {
		year := r.URL.Query().Get("y")
		month := r.URL.Query().Get("m")
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	}
}

// AdminReservationsCalendar shows the admin reservations calendar page.
func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {

	now := time.Now()
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	if year != "" && month != "" {
		year, err := strconv.Atoi(year)
		if err != nil {
			m.App.ErrorLog.Println(err)
			m.App.Session.Put(r.Context(), "error", "Error processing params, make sure you use the correct format")
			http.Redirect(w, r, "admin/reservations-calendar", http.StatusSeeOther)
			return
		}
		month, err := strconv.Atoi(month)
		if err != nil {
			m.App.ErrorLog.Println(err)
			m.App.Session.Put(r.Context(), "error", "Error processing params, make sure you use the correct format")
			http.Redirect(w, r, "/admin/reservations-calendar", http.StatusSeeOther)
			return
		}

		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")
	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := map[string]string{
		"next-month":      nextMonth,
		"next-month-year": nextMonthYear,
		"last-month":      lastMonth,
		"last-month-year": lastMonthYear,
		"this-month":      now.Format("01"),
		"this-month-year": now.Format("2006"),
	}

	// get the first and last days of the month
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := map[string]int{
		"days-in-month": lastOfMonth.Day(),
	}

	rooms, err := m.DB.GetAllRooms()
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error retrieving rooms")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"now":   now,
		"rooms": rooms,
	}

	for _, room := range rooms {
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for date := firstOfMonth; !date.After(lastOfMonth); date = date.AddDate(0, 0, 1) {
			reservationMap[date.Format("2006-01-2")] = 0
			blockMap[date.Format("2006-01-2")] = 0
		}

		restrictions, err := m.DB.GetRestrictionsForRoomByDate(room.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			m.App.ErrorLog.Println(err)
			m.App.Session.Put(r.Context(), "error", "Error retrieving restrictions")
			http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
			return
		}

		for _, restriction := range restrictions {
			if restriction.ReservationID > 0 {
				// it's a reservation
				for date := restriction.StartDate; !date.After(restriction.EndDate); date = date.AddDate(0, 0, 1) {
					reservationMap[date.Format("2006-01-2")] = restriction.ReservationID
				}
			} else {
				// it's a block
				blockMap[restriction.StartDate.Format("2006-01-2")] = restriction.ID
			}
		}

		data[fmt.Sprintf("reservation-map-%d", room.ID)] = reservationMap
		data[fmt.Sprintf("block-map-%d", room.ID)] = blockMap

		m.App.Session.Put(r.Context(), fmt.Sprintf("block-map-%d", room.ID), blockMap)
	}

	render.Template(w, r, "admin-reservations-calendar.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
		IntMap:    intMap,
	})
}

// AdminPostReservationsCalendar handles the updating of the reservations calendar
func (m *Repository) AdminPostReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error retrieving form information")
		http.Redirect(w, r, "/admin/reservations-calendar", http.StatusSeeOther)
		return
	}

	year, err := strconv.Atoi(r.Form.Get("y"))
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error retrieving form information")
		http.Redirect(w, r, "/admin/reservations-calendar", http.StatusSeeOther)
		return
	}
	month, err := strconv.Atoi(r.Form.Get("m"))
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error retrieving form information")
		http.Redirect(w, r, "/admin/reservations-calendar", http.StatusSeeOther)
		return
	}

	//process blocks
	rooms, err := m.DB.GetAllRooms()
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Error retrieving form information")
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
		return
	}

	form := forms.New(r.PostForm)

	// remove blocks
	for _, room := range rooms {
		currentMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block-map-%d", room.ID)).(map[string]int)
		for key, value := range currentMap {
			if val, ok := currentMap[key]; ok {
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove-block-%d-%s", room.ID, key)) {
						err := m.DB.DeleteBlockByID(value)
						if err != nil {
							m.App.ErrorLog.Println(err)
							m.App.Session.Put(r.Context(), "error", "Error deleting block")
							http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
							return
						}
					}
				}
			}
		}
	}

	for name := range r.PostForm {
		if strings.HasPrefix(name, "add-block") {
			splitName := strings.Split(name, "-")
			roomID, err := strconv.Atoi(splitName[2])
			if err != nil {
				m.App.ErrorLog.Println(err)
				m.App.Session.Put(r.Context(), "error", "Error parsing room id")
				http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusTemporaryRedirect)
				return
			}
			timeLayout := "2006-01-2"
			date, err := time.Parse(timeLayout, strings.Join(splitName[3:], "-"))
			if err != nil {
				m.App.ErrorLog.Println(err)
				m.App.Session.Put(r.Context(), "error", "Error parsing time")
				http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusTemporaryRedirect)
				return
			}
			err = m.DB.InsertBlockForRoom(roomID, date)
			if err != nil {
				m.App.ErrorLog.Println(err)
				m.App.Session.Put(r.Context(), "error", "Error inserting block")
				http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusTemporaryRedirect)
				return
			}
		}
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
}
