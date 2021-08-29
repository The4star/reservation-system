package render

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/justinas/nosurf"
	"github.com/the4star/reservation-system/internal/config"
	"github.com/the4star/reservation-system/internal/models"
)

var functions = template.FuncMap{}
var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(data *models.TemplateData, r *http.Request) *models.TemplateData {
	data.Flash = app.Session.PopString(r.Context(), "flash")
	data.Error = app.Session.PopString(r.Context(), "error")
	data.Warning = app.Session.PopString(r.Context(), "warning")
	data.CSRFToken = nosurf.Token((r))
	return data
}

// renders html templates
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data *models.TemplateData) {
	var templateCache map[string]*template.Template
	var err error
	if app.UseCache {
		//get the template cache from the app config
		templateCache = app.TemplateCache
	} else {
		// rebuid cache on every request
		templateCache, err = CreateTemplateCache()
		if err != nil {
			log.Fatal("unable to create fresh template cache")
		}
	}

	template, ok := templateCache[tmpl]
	if !ok {
		log.Fatal("unable to find template from cache")
	}

	buf := new(bytes.Buffer)
	data = AddDefaultData(data, r)
	err = template.Execute(buf, data)
	if err != nil {
		fmt.Println("Error writing template to browser:", err)
	}

	buf.WriteTo(w)
}

// creates a template cache as a map.
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	childPages, err := filepath.Glob("./templates/*/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	pages = append(pages, childPages...)

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}

	return myCache, nil
}
