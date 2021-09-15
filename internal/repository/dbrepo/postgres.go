package dbrepo

import (
	"context"
	"time"

	"github.com/the4star/reservation-system/internal/models"
)

func (pr *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the Database.
func (pr *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `insert into reservations(first_name, last_name, email,
		phone, start_date, end_date, room_id, created_at, updated_at)
		values($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id
	`
	err := pr.DB.QueryRowContext(
		ctx,
		stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

//InsertRoomRestriction inserts a room into the database
func (pr *postgresDBRepo) InsertRoomRestriction(rr models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id,
		created_at, updated_at, restriction_id) 
		values($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := pr.DB.ExecContext(
		ctx,
		stmt,
		rr.StartDate,
		rr.EndDate,
		rr.RoomID,
		rr.ReservationID,
		time.Now(),
		time.Now(),
		rr.RestrictionID,
	)

	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDates returns true if availability exists and false if no availability exists
func (pr *postgresDBRepo) SearchAvailabilityByDates(roomID int, startDate, endDate time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select count(id)
		from room restrictions
		where room_id = $1
		and $2 < end_date $3 > start_date 
	`

	var numRows int

	row, err := pr.DB.QueryContext(ctx, query, roomID, startDate, endDate)
	if err != nil {
		return false, err
	}
	err = row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}
