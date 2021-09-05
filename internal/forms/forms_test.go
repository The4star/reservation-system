package forms

import (
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	request := httptest.NewRequest("POST", "/whatever", nil)
	form := New(request.PostForm)
	formType := reflect.ValueOf(form).Elem().Type()
	isForm := formType == reflect.TypeOf(Form{})
	if !isForm {
		t.Errorf("Expected a form but got %T", formType)
	}
}

func TestValid(t *testing.T) {
	request := httptest.NewRequest("POST", "/whatever", nil)
	form := New(request.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("Got invalid when should have been valid")
	}
}

func TestRequired(t *testing.T) {
	request := httptest.NewRequest("POST", "/whatever", nil)
	form := New(request.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form is valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	request = httptest.NewRequest("POST", "/whatever", nil)
	request.PostForm = postedData
	form = New(request.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form showing invalid when should be valid")
	}
}

func TestHas(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("a", "")
	form := New(postedData)
	form.Has("a")
	if form.Valid() {
		t.Error("Form is valid when should be invalid")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)
	form.Has("a")
	if !form.Valid() {
		t.Error("form is invalid when should be valid")
	}
}

func TestMinLength(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("a", "a")
	form := New(postedData)
	form.MinLength("a", 3)
	if form.Valid() {
		t.Error("Form is valid when should be invalid")
	}

	postedData = url.Values{}
	postedData.Add("a", "abc")
	form = New(postedData)
	form.MinLength("a", 3)
	if !form.Valid() {
		t.Error("form is invalid when should be valid")
	}
}

func TestIsEmail(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("a", "a")
	form := New(postedData)
	form.IsEmail("a")
	if form.Valid() {
		t.Error("Form is valid when should be invalid")
	}

	postedData = url.Values{}
	postedData.Add("a", "clinton@test.com")
	form = New(postedData)
	form.IsEmail("a")
	if !form.Valid() {
		t.Error("form is invalid when should be valid")
	}
}

func TestErrorsGet(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("a", "a")
	form := New(postedData)
	form.IsEmail("a")
	if form.Errors.Get("a") == "" {
		t.Error("Error is empty when should not be")
	}

	if form.Errors.Get("b") != "" {
		t.Error("got an error when expected empty string")
	}

}
