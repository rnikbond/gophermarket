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

	query := `CREATE TABLE IF NOT EXISTS users (
                           id SERIAL PRIMARY KEY,
		             username CHARACTER VARYING(50),
		        password_hash CHARACTER VARYING(64) );`

	_, err := db.Exec(query)
	return err
}
