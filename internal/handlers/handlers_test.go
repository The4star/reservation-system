package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

type test struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
	postType           string
}

var theTests = []test{
	{"home", "/", "GET", []postData{}, http.StatusOK, "none"},
	{"about", "/about", "GET", []postData{}, http.StatusOK, "none"},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK, "none"},
	{"standard suite", "/rooms/standard-suite", "GET", []postData{}, http.StatusOK, "none"},
	{"deluxe suite", "/rooms/deluxe-suite", "GET", []postData{}, http.StatusOK, "none"},
	{"book", "/book", "GET", []postData{}, http.StatusOK, "none"},
	{"reservation summary", "/reservation-summary", "GET", []postData{}, http.StatusOK, "none"},
	{"availablity", "/availability", "GET", []postData{}, http.StatusOK, "none"},
	{"check availability", "/availability", "POST", []postData{
		{key: "start-date", value: "2022-01-01"},
		{key: "end-date", value: "2022-01-02"},
	}, http.StatusOK, "form"},
	{"make a booking", "/book", "POST", []postData{
		{key: "first-name", value: "Clinton"},
		{key: "last-name", value: "Forster"},
		{key: "phone", value: "04123456789"},
		{key: "email", value: "clinton@test.com"},
	}, http.StatusOK, "form"},
	{"check room availability", "/room-availability", "POST", []postData{
		{key: "startDate", value: "2022-01-01"},
		{key: "endDate", value: "2022-01-02"},
		{key: "roomType", value: "Deluxe Suite"},
	}, http.StatusOK, "json"},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range theTests {
		if test.method == "GET" {
			// test get route
			resp, err := testServer.Client().Get(testServer.URL + test.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != test.expectedStatusCode {
				t.Errorf("expected %d for %s route but got %d", test.expectedStatusCode, test.name, resp.StatusCode)
			}
		} else {
			// test post route
			if test.postType == "form" {
				//form test
				values := url.Values{}
				for _, param := range test.params {
					values.Add(param.key, param.value)
				}
				resp, err := testServer.Client().PostForm(testServer.URL+test.url, values)
				if err != nil {
					t.Log(err)
					t.Fatal(err)
				}

				if resp.StatusCode != test.expectedStatusCode {
					t.Errorf("expected %d for %s route but got %d", test.expectedStatusCode, test.name, resp.StatusCode)
				}
			} else {
				// json test
				data := make(map[string]string)
				for _, param := range test.params {
					data[param.key] = param.value
				}
				jsonData, err := json.Marshal(data)
				if err != nil {
					t.Log(err)
					t.Fatal(err)
				}
				body := bytes.NewBuffer(jsonData)

				resp, err := testServer.Client().Post(testServer.URL+test.url, "application/json", body)
				if err != nil {
					t.Log(err)
					t.Fatal(err)
				}
				if resp.StatusCode != test.expectedStatusCode {
					t.Errorf("expected %d for %s route but got %d", test.expectedStatusCode, test.name, resp.StatusCode)
				}
			}
		}
	}
}
