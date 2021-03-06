package repository

import (
	"time"

	"github.com/the4star/reservation-system/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	GetRestrictionsForRoomByDate(roomId int, start, end time.Time) ([]models.RoomRestriction, error)
	InsertRoomRestriction(rr models.RoomRestriction) error
	InsertBlockForRoom(id int, startDate time.Time) error
	DeleteBlockByID(id int) error
	SearchAvailabilityByDatesByRoomID(roomID int, startDate, endDate time.Time) (bool, error)

	SearchAvailabilityForAllRooms(startDate, endDate time.Time) ([]models.Room, error)
	GetAllRooms() ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)

	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, password string) (int, string, error)

	GetAllReservations() ([]models.Reservation, error)
	GetAllNewReservations() ([]models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error)
	UpdateReservation(r models.Reservation) error
	DeleteReservation(id int) error
	UpdateProcessedForReservation(id int, processed bool) error
}
