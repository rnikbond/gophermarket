package postgres

import (
	"time"

	market "gophermarket/pkg"
)

func (pg Postgres) CreateOrder(number int64, username string) error {

	queryUserID := `SELECT id FROM users WHERE username = $1`
	var userID int64
	row := pg.db.QueryRow(queryUserID, username)
	if err := row.Scan(&userID); err != nil {
		return market.ErrUserNotFound
	}

	queryUserOrder := `SELECT user_id FROM orders WHERE number = $1`
	var orderUserID int64
	row = pg.db.QueryRow(queryUserOrder, number)

	if err := row.Scan(&userID); err == nil {
		// Номер заказа уже существует. Оставлось проверить, кто его создал

		if orderUserID == userID {
			// Этот пользователь уже делал этот заказ
			return nil
		}

		// Кто-то другой уже делал такой заказ
		return market.ErrOrderAlreadyExists
	}

	// Если дошли сюда - значит такого заказа еще не было - создаем

	query := `INSERT INTO orders (user_id, number,status,created_at) VALUES($1,$2,$3,$4)`
	_, err := pg.db.Exec(query, userID, number, "NEW", time.Now().Format(time.RFC3339))

	return err
}
