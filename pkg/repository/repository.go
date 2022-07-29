package repository

import (
	market "gophermarket/pkg"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	Login(user market.User) error
}

type Registration interface {
	CreateUser(user market.User) error
}

type Repository struct {
	Authorization
	Registration
}

func NewRepository(db *sqlx.DB) *Repository {

	return &Repository{
		Authorization: NewLoginService(db),
		Registration:  NewRegisterPostgres(db),
	}
}
