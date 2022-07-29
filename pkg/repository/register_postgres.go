package repository

import (
	"fmt"

	market "gophermarket/pkg"

	"github.com/jmoiron/sqlx"
)

type RegisterPostgres struct {
	db *sqlx.DB
}

func NewRegisterPostgres(db *sqlx.DB) *RegisterPostgres {
	return &RegisterPostgres{db: db}
}

func (r *RegisterPostgres) CreateUser(user market.User) error {

	querySelect := fmt.Sprintf("SELECT username, password_hash FROM %s", usersTable)
	row := r.db.QueryRow(querySelect)

	var userDB market.User

	if err := row.Scan(&userDB.Username, &userDB.Password); err == nil {
		return market.ErrUserAlreadyExists
	}

	queryInsert := fmt.Sprintf("INSERT INTO %s (username, password_hash) VALUES ($1, $2)", usersTable)

	if _, err := r.db.Exec(queryInsert, user.Username, user.Password); err != nil {
		return err
	}
	return nil
}
