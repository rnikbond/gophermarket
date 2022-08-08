package postgres

import (
	"log"
	"time"

	"gophermarket/internal/repository"
	market "gophermarket/pkg"
	pkgOrder "gophermarket/pkg/order"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

	_, err := pg.db.Exec(queryCreateOrder, userID, number, pkgOrder.StatusNew, time.Now().Format(time.RFC3339))
	return err
}

// GetByStatuses - Получение заказов с запрашиваемым статусом
func (pg Order) GetByStatuses(statuses []string) (map[int64]string, error) {

	rows, err := pg.db.Query(queryOrdersByStatuses, pq.Array(&statuses))
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			logrus.Printf("error close rows: %v\n", err)
		}
	}()

	orders := make(map[int64]string)

	for rows.Next() {
		var orderNum int64
		var status string

		if err := rows.Scan(&orderNum, &status); err != nil {
			return nil, err
		}

		orders[orderNum] = status
	}

	return orders, nil
}

// SetStatus - Изменение статуса заказа
func (pg Order) SetStatus(order int64, status string) error {

	_, err := pg.db.Exec(queryUpdateOrder, status, order)
	return err
}

// SetAccrual - Изменение начислений по заказу
func (pg Order) SetAccrual(order int64, accrual float64) error {

	_, err := pg.db.Exec(queryUpdateAccrual, int64(accrual*100), order)
	return err
}

func (pg Order) Accruals(username string) (float64, error) {

	var userID int64
	row := pg.db.QueryRow(queryGetUserIDByName, username)
	if err := row.Scan(&userID); err != nil {
		return 0, market.ErrUserNotFound
	}

	row = pg.db.QueryRow(queryUserAccruals, userID)

	var accrualsUser int64
	if err := row.Scan(&accrualsUser); err != nil {
		// TODO Если у пользователя нет было, то Scan возвращает
		// error при записи NULL в int64
		log.Printf("error get accruals: %v\n", err)
		return 0, nil
	}

	return float64(accrualsUser / 100), nil
}

func (pg Order) Withdrawals(username string) (float64, error) {

	var userID int64
	row := pg.db.QueryRow(queryGetUserIDByName, username)
	if err := row.Scan(&userID); err != nil {
		return 0, market.ErrUserNotFound
	}

	return 0, nil
}
