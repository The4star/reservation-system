package repository

import (
	"time"

	"github.com/the4star/reservation-system/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(rr models.RoomRestriction) error
	SearchAvailabilityByDates(roomID int, startDate, endDate time.Time) (bool, error)
}
