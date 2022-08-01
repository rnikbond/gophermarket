package repository

import (
	market "gophermarket/pkg"

	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) SignIn(user market.User) error {

	row := r.db.QueryRow(queryGetUser, user.Username)
	var userDB market.User

	if err := row.Scan(&userDB.Username, &userDB.Password); err != nil {
		return market.ErrUserNotFound
	}

	if user.Username != userDB.Username || user.Password != userDB.Password {
		return market.ErrUserNotFound
	}

	return nil
}

func (r *AuthPostgres) SignUp(user market.User) error {

	row := r.db.QueryRow(queryGetUser, user.Username)

	var userDB market.User

	if err := row.Scan(&userDB.Username, &userDB.Password); err == nil {
		return market.ErrUserAlreadyExists
	}

	if _, err := r.db.Exec(queryCreateUse, user.Username, user.Password); err != nil {
		return err
	}
	return nil
}
