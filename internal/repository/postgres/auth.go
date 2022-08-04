package postgres

import (
	market "gophermarket/internal"
	"gophermarket/internal/repository"
	"gophermarket/pkg"

	"github.com/jmoiron/sqlx"
)

type Authorization struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) repository.Authorization {
	return &Authorization{
		db: db,
	}
}

// Create Создание нового пользователя
func (pg Authorization) Create(user market.User) error {

	if ok := pg.ExistsUsername(user.Username); ok {
		return pkg.ErrUserAlreadyExists
	}

	if _, err := pg.db.Exec(queryCreateUser, user.Username, user.Password); err != nil {
		return err
	}
	return nil
}

// Id Получение идентификатора пользователя
func (pg Authorization) Id(user market.User) (int64, error) {
	row := pg.db.QueryRow(queryGetUserID, user.Username, user.Password)

	var userID int64
	if err := row.Scan(&userID); err != nil {
		return 0, pkg.ErrUserNotFound
	}

	return userID, nil
}

// ExistsUsername Проверка существования имени пользователя
func (pg Authorization) ExistsUsername(username string) bool {
	row := pg.db.QueryRow(queryGetUserIDByName, username)

	var userID int64
	if err := row.Scan(&userID); err != nil {
		return false
	}

	return true
}
