package postgres

import (
	"gophermarket/pkg/repository"

	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	db *sqlx.DB
}

func (pg *Postgres) Finish() error {
	return pg.db.Close()
}

func NewPostgresRepository(dsn string) (repository.Repository, error) {

	db, err := pgDriver(dsn)
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	pg := Postgres{
		db: db,
	}

	return &pg, nil
}

func pgDriver(dsn string) (*sqlx.DB, error) {

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

	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
              id SERIAL PRIMARY KEY,
              username CHARACTER VARYING(50),
         password_hash CHARACTER VARYING(64)
         );`,

		`CREATE TABLE IF NOT EXISTS orders (
                 id SERIAL PRIMARY KEY,
            user_id INTEGER REFERENCES users (id),
             number BIGINT UNIQUE,
             status CHARACTER VARYING(50),
         created_at TIMESTAMP
         );`,
	}

	for _, query := range tables {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}
