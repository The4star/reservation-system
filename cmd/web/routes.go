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

	protectedRouteGroup := router.Group(nil)
	protectedRouteGroup.Use(Auth)

	//routes
	router.Get("/", handlers.Repo.Home)
	router.Get("/about", handlers.Repo.About)
	router.Get("/contact", handlers.Repo.Contact)
	router.Get("/rooms/standard-suite", handlers.Repo.StandardSuite)
	router.Get("/rooms/deluxe-suite", handlers.Repo.DeluxeSuite)
	router.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	router.Get("/book-room", handlers.Repo.BookRoom)

	//booking
	router.Get("/book", handlers.Repo.Book)
	router.Post("/book", handlers.Repo.PostBook)
	router.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	//user
	router.Get("/user/login", handlers.Repo.ShowLogin)
	router.Post("/user/login", handlers.Repo.PostLogin)
	router.Get("/user/logout", handlers.Repo.Logout)

	//protected
	protectedRouteGroup.Get("/admin/dashboard", handlers.Repo.AdminDashboard)
	protectedRouteGroup.Get("/admin/reservations-new", handlers.Repo.AdminNewReservations)
	protectedRouteGroup.Get("/admin/reservations-all", handlers.Repo.AdminAllReservations)
	protectedRouteGroup.Get("/admin/reservations-calendar", handlers.Repo.AdminReservationsCalendar)
	protectedRouteGroup.Post("/admin/reservations-calendar", handlers.Repo.AdminPostReservationsCalendar)
	protectedRouteGroup.Get("/admin/reservations/{src}/{id}", handlers.Repo.AdminShowReservation)
	protectedRouteGroup.Post("/admin/reservations/{src}/{id}", handlers.Repo.AdminPostUpdateReservation)
	protectedRouteGroup.Get("/admin/process/{src}/{id}", handlers.Repo.AdminProcessReservation)
	protectedRouteGroup.Get("/admin/delete/{src}/{id}", handlers.Repo.AdminDeleteReservation)
	//availability
	noSurfGroup.Get("/availability", handlers.Repo.Availability)
	noSurfGroup.Post("/availability", handlers.Repo.PostAvailability)
	router.Post("/room-availability", handlers.Repo.PostRoomAvailability)

	fileServer := http.FileServer(http.Dir("./static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return router
}
