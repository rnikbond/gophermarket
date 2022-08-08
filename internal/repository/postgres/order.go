package postgres

import (
	"strconv"
	"time"

	"gophermarket/internal/repository"
	market "gophermarket/pkg"
	"gophermarket/pkg/logpack"
	pkgOrder "gophermarket/pkg/order"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Order struct {
	logger *logpack.LogPack
	db     *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB, logger *logpack.LogPack) repository.Order {
	return &Order{
		db:     db,
		logger: logger,
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

	_, err := pg.db.Exec(queryCreateOrder, userID, number, pkgOrder.StatusNew, time.Now().Format("2006-01-02T15:04:05Z07:00"))
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
			pg.logger.Err.Printf("could not close rows: %s\n", err)
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

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}

	return orders, nil
}

// SetStatus - Изменение статуса заказа
func (pg Order) SetStatus(order int64, status string) error {

	_, err := pg.db.Exec(queryUpdateOrder, status, order)
	return err
}

func (pg Order) UserOrders(username string) ([]pkgOrder.InfoOrder, error) {

	var userID int64
	row := pg.db.QueryRow(queryGetUserIDByName, username)
	if err := row.Scan(&userID); err != nil {
		return nil, market.ErrUserNotFound
	}

	rows, err := pg.db.Query(queryUserOrders, userID)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			pg.logger.Err.Printf("could not close rows: %s\n", err)
		}
	}()

	var infoOrders []pkgOrder.InfoOrder

	for rows.Next() {
		var infoOrder pkgOrder.InfoOrder
		var orderNum int64

		errScan := rows.Scan(&orderNum, &infoOrder.Status, &infoOrder.Accrual, &infoOrder.UploadedAt)
		if errScan != nil {
			return nil, errScan
		}

		infoOrder.Order = strconv.FormatInt(orderNum, 10)
		infoOrder.Accrual = infoOrder.Accrual / 100

		infoOrders = append(infoOrders, infoOrder)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return infoOrders, nil
}
