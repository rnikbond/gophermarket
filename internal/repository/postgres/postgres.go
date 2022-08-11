package postgres

import (
	"gophermarket/internal/repository"
	"gophermarket/pkg/logpack"

	"github.com/jmoiron/sqlx"
)

func NewPostgresRepository(db *sqlx.DB, logger *logpack.LogPack) (*repository.Repository, error) {

	if err := createTables(db); err != nil {
		return nil, err
	}

	return &repository.Repository{
		Authorization: NewAuthPostgres(db, logger),
		Order:         NewOrderPostgres(db, logger),
		Loyalty:       NewLoyaltyPostgres(db, logger),
	}, nil
}

func createTables(db *sqlx.DB) error {

	// TODO Заменить на миграции

	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
                     id SERIAL PRIMARY KEY,
               username CHARACTER VARYING(50),
          password_hash CHARACTER VARYING(64)
        );`,

		`CREATE TABLE IF NOT EXISTS orders (
                    id SERIAL PRIMARY KEY,
               user_id INTEGER REFERENCES users (id),
                number CHARACTER VARYING(50),
           uploaded_at TIMESTAMPTZ,
                status CHARACTER VARYING(50),
               accrual INTEGER DEFAULT 0,
            withdrawal INTEGER DEFAULT 0
        );`,
	}

	for _, query := range tables {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}
