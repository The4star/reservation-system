package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/the4star/reservation-system/internal/config"
)

var app *config.AppConfig

// NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func IsAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "user_id")
}

//NiceDate returns date in nice format.
func NiceDate(t time.Time) string {
	return t.Format("02-01-2006")
}

// FormatDate formats a date in the specified format(f)
func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

// returns a slice of integers starting at 1 and going to count.
func Iterate(count int) []int {
	var i int
	var items []int
	for i = 0; i < count; i++ {
		items = append(items, i+1)
	}
	return items
}
