package postgres

import (
	"context"
	"strconv"
	"time"

	"gophermarket/internal/repository"
	market "gophermarket/pkg"
	"gophermarket/pkg/logpack"

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
func (pg Order) Create(ctx context.Context, number int64, username string, status string) error {

	var userID int64
	row := pg.db.QueryRowContext(ctx, queryGetUserIDByName, username)
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

	_, err := pg.db.Exec(queryCreateOrder, userID, number, status, time.Now().Format(time.RFC3339))
	return err
}

func (pg Order) CreateWithPayment(ctx context.Context, number int64, username string, sum float64) error {

	if err := pg.Create(ctx, number, username, market.StatusProcessed); err != nil {
		if err == market.ErrUserAlreadyOrderedIt {
			return market.ErrOrderAlreadyExists
		}

		return err
	}

	_, err := pg.db.Exec(queryChangeWithdrawals, int64(sum*100), number)
	return err
}

func (pg Order) GetByStatuses(ctx context.Context, statuses []string) (map[int64]string, error) {

	rows, err := pg.db.QueryxContext(ctx, queryOrdersByStatuses, pq.Array(&statuses))
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
func (pg Order) SetStatus(ctx context.Context, order int64, status string) error {

	_, err := pg.db.ExecContext(ctx, queryUpdateOrder, status, order)
	return err
}

func (pg Order) UserOrders(ctx context.Context, username string) ([]market.OrderInfo, error) {

	var userID int64
	row := pg.db.QueryRowContext(ctx, queryGetUserIDByName, username)
	if err := row.Scan(&userID); err != nil {
		return nil, market.ErrUserNotFound
	}

	rows, err := pg.db.QueryContext(ctx, queryUserOrders, userID)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			pg.logger.Err.Printf("could not close rows: %s\n", err)
		}
	}()

	var infoOrders []market.OrderInfo

	for rows.Next() {
		var infoOrder market.OrderInfo
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
