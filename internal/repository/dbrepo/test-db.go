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

func (pr *testDBRepo) GetRestrictionsForRoomByDate(roomId int, start, end time.Time) ([]models.RoomRestriction, error) {
	var restrictions []models.RoomRestriction
	return restrictions, nil
}

//InsertRoomRestriction inserts a room into the database
func (pr *testDBRepo) InsertRoomRestriction(rr models.RoomRestriction) error {
	if rr.RoomID == 1000 {
		return errors.New("failed to insert room restriction")
	}
	return nil
}

// InsertBlockForRoom blocks out a room for the admin.
func (pr *testDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {
	return nil
}

// DeleteBlockByID deletes a room restriction by id.
func (pr *testDBRepo) DeleteBlockByID(id int) error {
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

func (pr *testDBRepo) GetAllRooms() ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomByID gets a room by idß
func (pr *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room

	if id > 2 {
		return room, errors.New("invalid Room")
	}
	return room, nil
}

// GetUserByID returns a user by id
func (pr *testDBRepo) GetUserByID(id int) (models.User, error) {
	var user models.User
	return user, nil
}

func (pr *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

//Authenticate authenticates a user
func (pr *testDBRepo) Authenticate(email, password string) (int, string, error) {
	var id int
	var hashedPassword string
	return id, hashedPassword, nil
}

//GetAllReservations gets all reservations from db.
func (pr *testDBRepo) GetAllReservations() ([]models.Reservation, error) {
	var allReservations []models.Reservation
	return allReservations, nil
}

//GetAllNewReservations gets all new reservations from db.
func (pr *testDBRepo) GetAllNewReservations() ([]models.Reservation, error) {
	var allReservations []models.Reservation
	return allReservations, nil
}

// GetReservationByID returns one reservation from db
func (pr *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	var reservation models.Reservation
	return reservation, nil
}

func (pr *testDBRepo) UpdateReservation(r models.Reservation) error {
	return nil
}

func (pr *testDBRepo) DeleteReservation(id int) error {
	return nil
}

func (pr *testDBRepo) UpdateProcessedForReservation(id int, processed bool) error {
	return nil
}
