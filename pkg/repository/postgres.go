package repository

import (
	"fmt"

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

	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sqlx.DB) error {

	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
                            id SERIAL PRIMARY KEY,
		                    username CHARACTER VARYING(50),
		                    password_hash CHARACTER VARYING(100));`, usersTable)

	_, err := db.Exec(query)
	return err
}
