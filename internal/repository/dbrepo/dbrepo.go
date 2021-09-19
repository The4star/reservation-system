package dbrepo

import (
	"database/sql"

	"github.com/the4star/reservation-system/internal/config"
	"github.com/the4star/reservation-system/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(db *sql.DB, app *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: app,
		DB:  db,
	}
}
