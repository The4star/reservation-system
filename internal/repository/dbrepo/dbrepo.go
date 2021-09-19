package dbrepo

import (
	"database/sql"

	"github.com/the4star/reservation-system/internal/config"
	"github.com/the4star/reservation-system/internal/repository"
)

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}
type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewTestingRepo(app *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: app,
	}
}

func NewPostgresRepo(db *sql.DB, app *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: app,
		DB:  db,
	}
}
