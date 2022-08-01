package repository

import (
	market "gophermarket/pkg"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	SignIn(user market.User) error
	SignUp(user market.User) error
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {

	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
