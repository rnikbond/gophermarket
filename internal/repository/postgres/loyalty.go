package postgres

import (
	"context"
	"strconv"

	"gophermarket/internal/repository"
	"gophermarket/pkg"
	"gophermarket/pkg/logpack"

	"github.com/jmoiron/sqlx"
)

type Loyalty struct {
	logger *logpack.LogPack
	db     *sqlx.DB
}

func NewLoyaltyPostgres(db *sqlx.DB, logger *logpack.LogPack) repository.Loyalty {
	return &Loyalty{
		db:     db,
		logger: logger,
	}
}

func (l Loyalty) HowMatchAvailable(ctx context.Context, username string) (float64, error) {

	var userID int64
	row := l.db.QueryRowContext(ctx, queryGetUserIDByName, username)
	if err := row.Scan(&userID); err != nil {
		return 0, pkg.ErrUserNotFound
	}

	row = l.db.QueryRow(queryUserAccruals, userID)

	var accrualsUser *int64
	if err := row.Scan(&accrualsUser); err != nil {
		l.logger.Err.Printf("could not get accruals: %s\n", err)
		return 0, err
	}

	if accrualsUser == nil {
		return 0, nil
	}

	return float64(*accrualsUser) / 100, nil
}

func (l Loyalty) HowMatchUsed(ctx context.Context, username string) (float64, error) {

	var userID int64
	row := l.db.QueryRowContext(ctx, queryGetUserIDByName, username)
	if err := row.Scan(&userID); err != nil {
		return 0, pkg.ErrUserNotFound
	}

	row = l.db.QueryRow(queryUserWithdrawals, userID)

	var withdrawnUser *int64
	if err := row.Scan(&withdrawnUser); err != nil {
		l.logger.Err.Printf("could not get withdrawn: %s\n", err)
		return 0, err
	}

	if withdrawnUser == nil {
		return 0, nil
	}

	return float64(*withdrawnUser) / 100, nil
}

// SetAccrual - Изменение начислений по заказу
func (l Loyalty) SetAccrual(ctx context.Context, order int64, accrual float64) error {

	accrualRound := int64(accrual * 100)
	_, err := l.db.ExecContext(ctx, queryUpdateAccrual, accrualRound, order)
	return err
}

// Payments - Получение информации о списании
func (l Loyalty) Payments(ctx context.Context, username string) ([]pkg.PaymentInfo, error) {

	var userID int64
	row := l.db.QueryRowContext(ctx, queryGetUserIDByName, username)
	if err := row.Scan(&userID); err != nil {
		return nil, pkg.ErrUserNotFound
	}

	rows, err := l.db.QueryContext(ctx, queryWithdrawalsInfo, userID)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			l.logger.Err.Printf("could not close rows: %s\n", err)
		}
	}()

	var paymentsInfo []pkg.PaymentInfo

	for rows.Next() {

		var orderNum int64
		var withdrawal int64

		var writeOff pkg.PaymentInfo

		if err := rows.Scan(&orderNum, &withdrawal, &writeOff.UploadedAt); err != nil {
			return nil, err
		}

		writeOff.OrderNum = strconv.FormatInt(orderNum, 10)
		writeOff.Sum = float64(withdrawal) / 100

		paymentsInfo = append(paymentsInfo, writeOff)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return paymentsInfo, nil
}
