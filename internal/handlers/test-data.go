package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type postData struct {
	key   string
	value string
}

type getTest struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}

type loginTest struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}

type updateResTest struct {
	name               string
	userFirstName      string
	userLastName       string
	userEmail          string
	userPhone          string
	src                string
	year               string
	month              string
	resID              int
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}

type processDeleteTest struct {
	name               string
	src                string
	resID              string
	year               string
	month              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}

type postCalendarTest struct {
	name                 string
	postedData           url.Values
	expectedResponseCode int
	blocks               int
	reservations         int
}

var getTests = []getTest{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"availability", "/availability", "GET", http.StatusOK},
	{"standard suite", "/rooms/standard-suite", "GET", http.StatusOK},
	{"deluxe suite", "/rooms/deluxe-suite", "GET", http.StatusOK},
	{"login", "/user/login", "GET", http.StatusOK},
	{"logout", "/user/logout", "GET", http.StatusOK},
	{"dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"new reservations", "/admin/reservations-new", "GET", http.StatusOK},
	{"all reservations", "/admin/reservations-all", "GET", http.StatusOK},
	{"show reservation", "/admin/reservations/new/1/show", "GET", http.StatusOK},
	{"calendar", "/admin/reservations-calendar", "GET", http.StatusOK},
	{"calendar with params", "/admin/reservations-calendar?y=2023&m=1", "GET", http.StatusOK},
	{"non existent", "/dog/toy", "GET", http.StatusNotFound},
}

var loginTests = []loginTest{
	{
		"valid-credentials",
		"me@here.com",
		http.StatusSeeOther,
		"",
		"/admin/dashboard",
	},
	{
		"invalid-credentials",
		"invalid@here.com",
		http.StatusOK,
		`action="/user/login"`,
		"",
	},
	{
		"invalid-data",
		"j",
		http.StatusOK,
		`action="/user/login"`,
		"",
	},
}

var updateReservationTests = []updateResTest{
	{
		"valid data",
		"John",
		"Smith",
		"me@here.com",
		"1234567",
		"new",
		"",
		"",
		1,
		http.StatusSeeOther,
		"",
		"",
	},
	{
		"from calendar",
		"John",
		"Smith",
		"me@here.com",
		"1234567",
		"cal",
		"2050",
		"1",
		1,
		http.StatusSeeOther,
		"",
		"",
	},
	{
		"invalid data",
		"John",
		"Smith",
		"j",
		"1234567",
		"new",
		"",
		"",
		1,
		http.StatusOK,
		"",
		"",
	},
	{
		"invalid reservation",
		"John",
		"Smith",
		"j",
		"1234567",
		"new",
		"",
		"",
		1000,
		http.StatusSeeOther,
		"",
		"",
	},
}

var procesDeleteTests = []processDeleteTest{
	{
		name:               "New as source",
		src:                "new",
		resID:              "1",
		year:               "",
		month:              "",
		expectedStatusCode: http.StatusSeeOther,
		expectedHTML:       "",
		expectedLocation:   "/admin/reservations-new",
	},
	{
		name:               "Cal as source",
		src:                "cal",
		resID:              "1",
		year:               "2050",
		month:              "1",
		expectedStatusCode: http.StatusSeeOther,
		expectedHTML:       "",
		expectedLocation:   "/admin/reservations-calendar?y=2050&m=1",
	},
	{
		name:               "invalid res id",
		src:                "cal",
		resID:              "1000",
		year:               "2050",
		month:              "1",
		expectedStatusCode: http.StatusSeeOther,
		expectedHTML:       "",
		expectedLocation:   "/admin/reservations/cal/1000",
	},
}

var adminPostReservationCalendarTests = []postCalendarTest{
	{
		name: "cal",
		postedData: url.Values{
			"y": {time.Now().Format("2006")},
			"m": {time.Now().Format("01")},
			fmt.Sprintf("add_block_1_%s", time.Now().AddDate(0, 0, 2).Format("2006-01-2")): {"1"},
		},
		expectedResponseCode: http.StatusSeeOther,
	},
	{
		name:                 "cal-blocks",
		postedData:           url.Values{},
		expectedResponseCode: http.StatusSeeOther,
		blocks:               1,
	},
	{
		name:                 "cal-res",
		postedData:           url.Values{},
		expectedResponseCode: http.StatusSeeOther,
		reservations:         1,
	},
}
