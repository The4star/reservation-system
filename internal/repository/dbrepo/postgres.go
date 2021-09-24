package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/the4star/reservation-system/internal/models"
	"golang.org/x/crypto/bcrypt"
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

//GetRestrictionsForRoomByDate returns restrictions for a room by date range.
func (pr *postgresDBRepo) GetRestrictionsForRoomByDate(roomId int, start, end time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restrictions []models.RoomRestriction

	query := `
	select id, coalesce(reservation_id, 0), restriction_id, room_id, start_date, end_date
	from room_restrictions where $1 < end_date and $2 > start_date
	and room_id = $3
	`

	rows, err := pr.DB.QueryContext(ctx, query, start, end, roomId)
	if err != nil {
		return restrictions, err
	}
	defer rows.Close()

	for rows.Next() {
		var restriction models.RoomRestriction
		err := rows.Scan(
			&restriction.ID,
			&restriction.ReservationID,
			&restriction.RestrictionID,
			&restriction.RoomID,
			&restriction.StartDate,
			&restriction.EndDate,
		)
		if err != nil {
			return restrictions, nil
		}
		restrictions = append(restrictions, restriction)
	}

	if err = rows.Err(); err != nil {
		return restrictions, err
	}

	return restrictions, nil
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

// SearchAvailabilityByDatesByRoomID returns true if availability exists and false if no availability exists
func (pr *postgresDBRepo) SearchAvailabilityByDatesByRoomID(roomID int, startDate, endDate time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select count(id)
		from room_restrictions
		where room_id = $1
		and $2 < end_date and $3 > start_date 
	`

	var numRows int

	row := pr.DB.QueryRowContext(ctx, query, roomID, startDate, endDate)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms if any for given date range
func (pr *postgresDBRepo) SearchAvailabilityForAllRooms(startDate, endDate time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select r.id, r.room_name
		from rooms r
		where r.id not in (select room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date) 
	`

	var availableRooms []models.Room

	rows, err := pr.DB.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return availableRooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return availableRooms, err
		}
		availableRooms = append(availableRooms, room)
	}

	if err = rows.Err(); err != nil {
		return availableRooms, err
	}

	return availableRooms, nil
}

// GetAllRooms gets all rooms from db
func (pr *postgresDBRepo) GetAllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var rooms []models.Room

	query := `select id, room_name, created_at, updated_at from rooms`

	rows, err := pr.DB.QueryContext(ctx, query)
	if err != nil {
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
			&room.CreatedAt,
			&room.UpdatedAt,
		)

		if err != nil {
			return rooms, err
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetRoomByID gets a room by id
func (pr *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room
	query := `select id, room_name, created_at, updated_at from rooms where id = $1`

	row := pr.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}

	return room, nil
}

// GetUserByID returns a user by id
func (pr *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at
		from users where id = $1
	`

	row := pr.DB.QueryRowContext(ctx, query, id)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (pr *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name = $1, last_name = $2, email = $3, access_level = $4 updated_at = $5
		where id = $6
	`

	_, err := pr.DB.ExecContext(
		ctx,
		query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
		u.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

//Authenticate authenticates a user
func (pr *postgresDBRepo) Authenticate(email, password string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "select id, password from users where email = $1"

	var id int
	var hashedPassword string

	row := pr.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

//GetAllReservations returns a slice of all reservations from db.
func (pr *postgresDBRepo) GetAllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var allReservations []models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
		r.end_date, r.room_id, r.created_at, r.updated_at,
		rm.id, rm.room_name
		from reservations r
		inner join rooms rm on (rm.id = r.room_id)
		order by r.start_date asc
	`

	rows, err := pr.DB.QueryContext(ctx, query)
	if err != nil {
		return allReservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var res models.Reservation
		err := rows.Scan(
			&res.ID,
			&res.FirstName,
			&res.LastName,
			&res.Email,
			&res.Phone,
			&res.StartDate,
			&res.EndDate,
			&res.RoomID,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.Room.ID,
			&res.Room.RoomName,
		)
		if err != nil {
			return allReservations, err
		}
		allReservations = append(allReservations, res)
	}

	if err = rows.Err(); err != nil {
		return allReservations, err
	}

	return allReservations, nil
}

//GetAllNewReservations returns a slice of all new reservations from db.
func (pr *postgresDBRepo) GetAllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var allReservations []models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
		r.end_date, r.room_id, r.created_at, r.updated_at, r.processed,
		rm.id, rm.room_name
		from reservations r
		inner join rooms rm on (rm.id = r.room_id)
		where processed = false
		order by r.start_date asc
	`

	rows, err := pr.DB.QueryContext(ctx, query)
	if err != nil {
		return allReservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var res models.Reservation
		err := rows.Scan(
			&res.ID,
			&res.FirstName,
			&res.LastName,
			&res.Email,
			&res.Phone,
			&res.StartDate,
			&res.EndDate,
			&res.RoomID,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.Processed,
			&res.Room.ID,
			&res.Room.RoomName,
		)
		if err != nil {
			return allReservations, err
		}
		allReservations = append(allReservations, res)
	}

	if err = rows.Err(); err != nil {
		return allReservations, err
	}

	return allReservations, nil
}

// GetReservationByID returns one reservation from db
func (pr *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.Reservation

	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date,
		r.room_id, r.created_at, r.updated_at, r.processed,
		rm.id, rm.room_name
		from reservations r
		inner join rooms rm on (r.room_id = rm.id)
		where r.id = $1
	`

	row := pr.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.RoomName,
	)
	if err != nil {
		return res, err
	}

	return res, nil
}

// DeleteReservation updates a reservation in the db
func (pr *postgresDBRepo) UpdateReservation(r models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set first_name = $1, last_name = $2, email = $3, phone = $4, updated_at = $5
		where id = $6
	`

	_, err := pr.DB.ExecContext(
		ctx,
		query,
		r.FirstName,
		r.LastName,
		r.Email,
		r.Phone,
		time.Now(),
		r.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteReservation deletes a reservation based on an id.
func (pr *postgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from reservations where id = $1`

	_, err := pr.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateProcessedForReservation updates processed for a reservation by id
func (pr *postgresDBRepo) UpdateProcessedForReservation(id int, processed bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set processed = $1 where id = $2`

	_, err := pr.DB.ExecContext(ctx, query, processed, id)
	if err != nil {
		return err
	}

	return nil
}
