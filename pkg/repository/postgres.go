package repository

import (
	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "users"
)

func NewPostgresDB(dsn string) (*sqlx.DB, error) {

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
