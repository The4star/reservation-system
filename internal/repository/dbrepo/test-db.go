package dbrepo

import (
	"errors"
	"time"

	"github.com/the4star/reservation-system/internal/models"
)

func (pr *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the Database.
func (pr *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID == 2 {
		return 0, errors.New("failed to insert reservation")
	}
	return 1, nil
}

//InsertRoomRestriction inserts a room into the database
func (pr *testDBRepo) InsertRoomRestriction(rr models.RoomRestriction) error {
	if rr.RoomID == 1000 {
		return errors.New("failed to insert room restriction")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists and false if no availability exists
func (pr *testDBRepo) SearchAvailabilityByDatesByRoomID(roomID int, startDate, endDate time.Time) (bool, error) {
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms if any for given date range
func (pr *testDBRepo) SearchAvailabilityForAllRooms(startDate, endDate time.Time) ([]models.Room, error) {
	var availableRooms []models.Room

	availableRooms = append(availableRooms, models.Room{
		ID:        1,
		RoomName:  "Standard Suite",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	})
	return availableRooms, nil
}

// GetRoomByID gets a room by idß
func (pr *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room

	if id > 2 {
		return room, errors.New("invalid Room")
	}
	return room, nil
}
