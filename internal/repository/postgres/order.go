package postgres

import (
	"time"

	"gophermarket/internal/repository"
	market "gophermarket/pkg"

	"github.com/jmoiron/sqlx"
)

type Order struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) repository.Order {
	return &Order{
		db: db,
	}
}

func (pg Order) Create(number int64, username string) error {

	var userID int64
	row := pg.db.QueryRow(queryGetUserIDByName, username)
	if err := row.Scan(&userID); err != nil {
		return market.ErrUserNotFound
	}

	var orderUserID int64
	row = pg.db.QueryRow(queryOrderUserID, number)
	if err := row.Scan(&orderUserID); err == nil {
		// Номер заказа уже существует. Осталось проверить, кто его создал

		if orderUserID == userID {
			// Этот пользователь уже делал этот заказ
			return market.ErrUserAlreadyOrderedIt
		}

		// Кто-то другой уже делал заказ с таким номером
		return market.ErrOrderAlreadyExists
	}

	// Если дошли сюда - значит такого заказа еще не было - создаем

	_, err := pg.db.Exec(queryCreateOrder, userID, number, "NEW", time.Now().Format(time.RFC3339))
	return err
}
