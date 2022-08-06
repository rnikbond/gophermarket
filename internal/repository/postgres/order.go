package postgres

import (
	"time"

	"gophermarket/internal/repository"
	market "gophermarket/pkg"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Order struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) repository.Order {
	return &Order{
		db: db,
	}
}

// Create - Создание нового заказа
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

	_, err := pg.db.Exec(queryCreateOrder, userID, number, "PROCESSING", time.Now().Format(time.RFC3339))
	return err
}

// GetByStatus - Получение заказов с запрашиваемым статусом
func (pg Order) GetByStatus(status string) ([]int64, error) {

	var orders []int64

	rows, err := pg.db.Query(queryOrdersByStatus, status)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			logrus.Printf("error close rows: %v\n", err)
		}
	}()

	for rows.Next() {
		var order int64
		if err := rows.Scan(&order); err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

// SetStatus - Изменение статуса заказа
func (pg Order) SetStatus(order int64, status string) error {

	_, err := pg.db.Exec(queryUpdateOrder, status, order)
	return err
}
