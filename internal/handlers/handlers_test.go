package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/the4star/reservation-system/internal/models"
)

func TestGetHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range getTests {
		// test get route
		fmt.Println(test.name)
		resp, err := testServer.Client().Get(testServer.URL + test.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != test.expectedStatusCode {
			t.Errorf("expected %d for %s route but got %d", test.expectedStatusCode, test.name, resp.StatusCode)
		}

		// test post route
		// if test.postType == "form" {
		// 	//form test
		// 	values := url.Values{}
		// 	for _, param := range test.params {
		// 		values.Add(param.key, param.value)
		// 	}
		// 	resp, err := testServer.Client().PostForm(testServer.URL+test.url, values)
		// 	if err != nil {
		// 		t.Log(err)
		// 		t.Fatal(err)
		// 	}

		// 	if resp.StatusCode != test.expectedStatusCode {
		// 		t.Errorf("expected %d for %s route but got %d", test.expectedStatusCode, test.name, resp.StatusCode)
		// 	}
		// } else {
		// 	// json test
		// 	data := make(map[string]string)
		// 	for _, param := range test.params {
		// 		data[param.key] = param.value
		// 	}
		// 	jsonData, err := json.Marshal(data)
		// 	if err != nil {
		// 		t.Log(err)
		// 		t.Fatal(err)
		// 	}
		// 	body := bytes.NewBuffer(jsonData)

		// 	resp, err := testServer.Client().Post(testServer.URL+test.url, "application/json", body)
		// 	if err != nil {
		// 		t.Log(err)
		// 		t.Fatal(err)
		// 	}
		// 	if resp.StatusCode != test.expectedStatusCode {
		// 		t.Errorf("expected %d for %s route but got %d", test.expectedStatusCode, test.name, resp.StatusCode)
		// 	}
		// }

	}
}

func TestBook(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			RoomName: "Standard Suite",
		},
	}

	req, _ := http.NewRequest("GET", "/book", nil)
	ctx := getCTX(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.Book)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusOK)
	}

	// test where reservation not in session
	req, _ = http.NewRequest("Get", "/book", nil)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusOK)
	}

	//test with non existent room
	req, _ = http.NewRequest("Get", "/book", nil)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusOK)
	}
}

func TestPostBook(t *testing.T) {
	timeLayout := "2006-01-02"
	startDate, err := time.Parse(timeLayout, "2050-02-01")
	if err != nil {
		fmt.Println(err)
	}
	endDate, err := time.Parse(timeLayout, "2050-02-03")
	if err != nil {
		fmt.Println(err)
	}

	reservation := models.Reservation{
		RoomID:    1,
		StartDate: startDate,
		EndDate:   endDate,
		Room: models.Room{
			RoomName: "Standard Suite",
		},
	}

	formData := []postData{
		{key: "room-id", value: "1"},
		{key: "first-name", value: "Clinton"},
		{key: "last-name", value: "Forster"},
		{key: "phone", value: "04123456789"},
		{key: "email", value: "clinton@test.com"},
	}

	values := url.Values{}
	for _, data := range formData {
		values.Add(data.key, data.value)
	}
	body := strings.NewReader(values.Encode())

	req, _ := http.NewRequest("POST", "/book", body)
	ctx := getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.PostBook)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for missing post body
	req, _ = http.NewRequest("POST", "/book", nil)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for missing session
	req, _ = http.NewRequest("POST", "/book", body)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for error when inserting room data
	body = strings.NewReader(values.Encode())
	req, _ = http.NewRequest("POST", "/book", body)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	reservation.RoomID = 2
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for error when inserting room restriction data
	body = strings.NewReader(values.Encode())
	req, _ = http.NewRequest("POST", "/book", body)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	reservation.RoomID = 1000

	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for missing form data
	values.Del("first-name")
	body = strings.NewReader(values.Encode())
	req, _ = http.NewRequest("POST", "/book", body)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Post Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusOK)
	}
}

func TestChooseRoom(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			RoomName: "Standard Suite",
		},
	}

	req, _ := http.NewRequest("Get", "/choose-room/1", nil)
	ctx := getCTX(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test without res in session
	req, _ = http.NewRequest("Get", "/choose-room/1", nil)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test without room in endpoint
	req, _ = http.NewRequest("Get", "/choose-room", nil)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestPostAvailability(t *testing.T) {
	formData := []postData{
		{key: "start-date", value: "2060-02-01"},
		{key: "end-date", value: "2060-02-04"},
	}

	values := url.Values{}
	for _, data := range formData {
		values.Add(data.key, data.value)
	}
	body := strings.NewReader(values.Encode())

	req, _ := http.NewRequest("POST", "/availability", body)
	ctx := getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Post Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusOK)
	}

	// test invalid start date
	values.Del("start-date")
	body = strings.NewReader(values.Encode())
	req, _ = http.NewRequest("POST", "/availability", body)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test invalid end date
	values.Del("end-date")
	values.Set("start-date", "2060-01-01")
	body = strings.NewReader(values.Encode())
	req, _ = http.NewRequest("POST", "/availability", body)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func TestPostRoomAvailability(t *testing.T) {
	data := RoomAvailabilityRequest{
		StartDate: "2060-01-01",
		EndDate:   "2060-01-02",
		RoomID:    "1",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Error("error marshalling json")
	}
	body := bytes.NewBuffer(jsonData)

	req, _ := http.NewRequest("POST", "/room-availability", body)
	ctx := getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostRoomAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Post Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusOK)
	}

	//test with invalid start date
	data = RoomAvailabilityRequest{
		StartDate: "",
		EndDate:   "2060-01-02",
		RoomID:    "1",
	}

	jsonData, err = json.Marshal(data)
	if err != nil {
		t.Error("error marshalling json")
	}
	body = bytes.NewBuffer(jsonData)
	req, _ = http.NewRequest("POST", "/room-availability", body)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Post Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusInternalServerError)
	}
}

func TestBookRoom(t *testing.T) {
	req, _ := http.NewRequest("GET", "/book-room?id=1&sd=2060-01-01&ed=2060-01-04", nil)
	ctx := getCTX(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test db failed
	req, _ = http.NewRequest("GET", "/book-room?id=100&sd=2060-01-01&ed=2060-01-04", nil)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestReservationSummary(t *testing.T) {
	timeLayout := "2006-01-02"
	startDate, err := time.Parse(timeLayout, "2050-02-01")
	if err != nil {
		fmt.Println(err)
	}
	endDate, err := time.Parse(timeLayout, "2050-02-03")
	if err != nil {
		fmt.Println(err)
	}

	reservation := models.Reservation{
		RoomID:    1,
		StartDate: startDate,
		EndDate:   endDate,
		Room: models.Room{
			RoomName: "Standard Suite",
		},
	}

	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := getCTX(req)
	req = req.WithContext(ctx)
	session.Put(req.Context(), "reservation", reservation)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Post Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusOK)
	}

	// test with no reservation in session
	req, _ = http.NewRequest("GET", "/reservation-summary", nil)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post Book handler returned wrong response code, got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestLogin(t *testing.T) {
	for _, lt := range loginTests {
		postedData := url.Values{}
		postedData.Add("email", lt.email)
		postedData.Add("password", "$Password123")

		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		ctx := getCTX(req)
		req = req.WithContext(ctx)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.PostLogin)
		handler.ServeHTTP(rr, req)

		if rr.Code != lt.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", lt.name, lt.expectedStatusCode, rr.Code)
		}

		if lt.expectedLocation != "" {
			actualLocation, _ := rr.Result().Location()
			if actualLocation.String() != lt.expectedLocation {
				t.Errorf("failed %s: expected location %s but got location %s", lt.name, lt.expectedLocation, actualLocation)
			}
		}

		if lt.expectedHTML != "" {
			html := rr.Body.String()
			if !strings.Contains(html, lt.expectedHTML) {
				t.Errorf("failed %s: expected html %s but got %s", lt.name, lt.expectedHTML, html)
			}
		}
	}
}

func TestPostUpdateReservation(t *testing.T) {
	for _, urt := range updateReservationTests {
		postedData := url.Values{}
		postedData.Add("first-name", urt.userFirstName)
		postedData.Add("last-name", urt.userLastName)
		postedData.Add("email", urt.userEmail)
		postedData.Add("phone", urt.userPhone)
		postedData.Add("year", urt.year)
		postedData.Add("month", urt.month)

		req, _ := http.NewRequest("POST", fmt.Sprintf("/admin/reservations/%s/%d", urt.src, urt.resID), strings.NewReader(postedData.Encode()))
		ctx := getCTX(req)
		req = req.WithContext(ctx)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminPostUpdateReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != urt.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", urt.name, urt.expectedStatusCode, rr.Code)
		}

		if urt.expectedLocation != "" {
			actualLocation, _ := rr.Result().Location()
			if actualLocation.String() != urt.expectedLocation {
				t.Errorf("failed %s: expected location %s but got location %s", urt.name, urt.expectedLocation, actualLocation)
			}
		}

		if urt.expectedHTML != "" {
			html := rr.Body.String()
			if !strings.Contains(html, urt.expectedHTML) {
				t.Errorf("failed %s: expected html %s but got %s", urt.name, urt.expectedHTML, html)
			}
		}
	}
}

func TestAdminProcessReservation(t *testing.T) {
	for _, prt := range procesDeleteTests {
		req, _ := http.NewRequest("POST", fmt.Sprintf("/admin/reservations/%s/%s/process?y=%s&m=%s", prt.src, prt.resID, prt.year, prt.month), nil)
		ctx := getCTX(req)
		req = req.WithContext(ctx)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminProcessReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != prt.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", prt.name, prt.expectedStatusCode, rr.Code)
		}

		if prt.expectedLocation != "" {
			actualLocation, _ := rr.Result().Location()
			if actualLocation.String() != prt.expectedLocation {
				t.Errorf("failed %s: expected location %s but got location %s", prt.name, prt.expectedLocation, actualLocation)
			}
		}

		if prt.expectedHTML != "" {
			html := rr.Body.String()
			if !strings.Contains(html, prt.expectedHTML) {
				t.Errorf("failed %s: expected html %s but got %s", prt.name, prt.expectedHTML, html)
			}
		}
	}
}

func TestAdminDeleteReservation(t *testing.T) {
	for _, prt := range procesDeleteTests {
		req, _ := http.NewRequest("POST", fmt.Sprintf("/admin/reservations/%s/%s/delete?y=%s&m=%s", prt.src, prt.resID, prt.year, prt.month), nil)
		ctx := getCTX(req)
		req = req.WithContext(ctx)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminDeleteReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != prt.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", prt.name, prt.expectedStatusCode, rr.Code)
		}

		if prt.expectedLocation != "" {
			actualLocation, _ := rr.Result().Location()
			if actualLocation.String() != prt.expectedLocation {
				t.Errorf("failed %s: expected location %s but got location %s", prt.name, prt.expectedLocation, actualLocation)
			}
		}

		if prt.expectedHTML != "" {
			html := rr.Body.String()
			if !strings.Contains(html, prt.expectedHTML) {
				t.Errorf("failed %s: expected html %s but got %s", prt.name, prt.expectedHTML, html)
			}
		}
	}
}

func TestPostReservationCalendar(t *testing.T) {
	for _, e := range adminPostReservationCalendarTests {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/admin/reservations-calendar", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/admin/reservations-calendar", nil)
		}
		ctx := getCTX(req)
		req = req.WithContext(ctx)

		now := time.Now()
		bm := make(map[string]int)
		rm := make(map[string]int)

		currentYear, currentMonth, _ := now.Date()
		currentLocation := now.Location()

		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			rm[d.Format("2006-01-2")] = 0
			bm[d.Format("2006-01-2")] = 0
		}

		if e.blocks > 0 {
			bm[firstOfMonth.Format("2006-01-2")] = e.blocks
		}

		if e.reservations > 0 {
			rm[lastOfMonth.Format("2006-01-2")] = e.reservations
		}

		session.Put(ctx, "block_map_1", bm)
		session.Put(ctx, "reservation_map_1", rm)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.AdminPostReservationsCalendar)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}

	}
}

func getCTX(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
