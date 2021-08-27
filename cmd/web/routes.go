package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/the4star/reservation-system/pkg/config"
	"github.com/the4star/reservation-system/pkg/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(NoSurf)
	router.Use(SessionLoad)

	//routes
	router.Get("/", handlers.Repo.Home)
	router.Get("/about", handlers.Repo.About)
	router.Get("/rooms/standard-suite", handlers.Repo.StandardSuite)
	router.Get("/rooms/deluxe-suite", handlers.Repo.DeluxeSuite)
	router.Get("/availability", handlers.Repo.Availability)
	router.Get("/book", handlers.Repo.Book)
	router.Get("/contact", handlers.Repo.Contact)

	fileServer := http.FileServer(http.Dir("./static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return router
}
