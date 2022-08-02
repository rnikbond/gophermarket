package postgres

import (
	market "gophermarket/pkg"
)

func (pg Postgres) SignUp(user market.User) error {

	row := pg.db.QueryRow(queryGetUser, user.Username)

	var userDB market.User

	if err := row.Scan(&userDB.Username, &userDB.Password); err == nil {
		return market.ErrUserAlreadyExists
	}

	if _, err := pg.db.Exec(queryCreateUse, user.Username, user.Password); err != nil {
		return err
	}
	return nil
}

func (pg Postgres) SignIn(user market.User) error {

	row := pg.db.QueryRow(queryGetUser, user.Username)
	var userDB market.User

	if err := row.Scan(&userDB.Username, &userDB.Password); err != nil {
		return market.ErrUserNotFound
	}

	if user.Username != userDB.Username || user.Password != userDB.Password {
		return market.ErrUserNotFound
	}

	return nil
}
