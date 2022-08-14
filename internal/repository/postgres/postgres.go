package postgres

import (
	"gophermarket/internal/repository"
	"gophermarket/pkg/logpack"

	"github.com/jmoiron/sqlx"
)

func NewPostgresRepository(db *sqlx.DB, logger *logpack.LogPack) (*repository.Repository, error) {

	return &repository.Repository{
		Authorization: NewAuthPostgres(db, logger),
		Order:         NewOrderPostgres(db, logger),
		Loyalty:       NewLoyaltyPostgres(db, logger),
	}, nil
}
