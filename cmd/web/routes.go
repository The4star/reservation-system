package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/the4star/reservation-system/internal/config"
	"github.com/the4star/reservation-system/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(SessionLoad)

	noSurfGroup := router.Group(nil)
	noSurfGroup.Use(NoSurf)

	//routes
	router.Get("/", handlers.Repo.Home)
	router.Get("/about", handlers.Repo.About)
	router.Get("/contact", handlers.Repo.Contact)
	router.Get("/rooms/standard-suite", handlers.Repo.StandardSuite)
	router.Get("/rooms/deluxe-suite", handlers.Repo.DeluxeSuite)
	router.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	router.Get("/book-room", handlers.Repo.BookRoom)

	router.Get("/book", handlers.Repo.Book)
	router.Post("/book", handlers.Repo.PostBook)
	router.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	router.Get("/user/login", handlers.Repo.ShowLogin)

	noSurfGroup.Get("/availability", handlers.Repo.Availability)
	noSurfGroup.Post("/availability", handlers.Repo.PostAvailability)
	router.Post("/room-availability", handlers.Repo.PostRoomAvailability)

	fileServer := http.FileServer(http.Dir("./static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return router
}
