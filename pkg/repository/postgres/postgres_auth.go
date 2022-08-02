package postgres

import (
	market "gophermarket/pkg"
)

func (pg Postgres) SignUp(user market.User) error {

	row := pg.db.QueryRow(queryCheckUsername, user.Username)

	var userID int64
	if err := row.Scan(&userID); err == nil {
		return market.ErrUserAlreadyExists
	}

	if _, err := pg.db.Exec(queryCreateUse, user.Username, user.Password); err != nil {
		return err
	}
	return nil
}

func (pg Postgres) SignIn(user market.User) error {

	row := pg.db.QueryRow(queryCheckUserAuth, user.Username, user.Password)

	var userID int64
	if err := row.Scan(&userID); err != nil {
		return market.ErrUserNotFound
	}

	return nil
}
