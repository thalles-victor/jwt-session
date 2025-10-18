package database

import (
	"jwt-session/src/utilities/config"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func Connect() (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", config.DATABASE_CONNECTION)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
