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

type postData struct {
	key   string
	value string
}

type test struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}

var getTests = []test{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"availability", "/availability", "GET", http.StatusOK},
	{"standard suite", "/rooms/standard-suite", "GET", http.StatusOK},
	{"deluxe suite", "/rooms/deluxe-suite", "GET", http.StatusOK},
	// {"reservation summary", "/reservation-summary", "GET", []postData{}, http.StatusOK, "none"},
	// {"availabilty", "/availability", "GET", []postData{}, http.StatusOK, "none"},
	// {"check availability", "/availability", "POST", []postData{
	// 	{key: "start-date", value: "2022-01-01"},
	// 	{key: "end-date", value: "2022-01-02"},
	// }, http.StatusOK, "form"},
	// {"check room availability", "/room-availability", "POST", []postData{
	// 	{key: "startDate", value: "2022-01-01"},
	// 	{key: "endDate", value: "2022-01-02"},
	// 	{key: "roomID", value: "2"},
	// }, http.StatusOK, "json"},
}

func TestGetHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range getTests {
		// test get route
		fmt.Println(test)
		resp, err := testServer.Client().Get(testServer.URL + test.url)
		if err != nil {
			t.Log("Here")
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

func getCTX(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
