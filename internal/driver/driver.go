package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

//DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

// ConnectSQl creates database tool for postgres
func ConnectSQl(dsn string) (*DB, error) {
	db, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(maxOpenDbConn)
	db.SetMaxIdleConns(maxIdleDbConn)
	db.SetConnMaxLifetime(maxDbLifetime)
	dbConn.SQL = db

	err = testDb(db)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

// NewDatabase creates a new instance of the database
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = testDb(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// testDb tests the database connection
func testDb(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}
