package postgres

import (
	"gophermarket/internal/repository"
	market "gophermarket/pkg"
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

func (l Loyalty) HowMatchAvailable(username string) (float64, error) {

	var userID int64
	row := l.db.QueryRow(queryGetUserIDByName, username)
	if err := row.Scan(&userID); err != nil {
		return 0, market.ErrUserNotFound
	}

	row = l.db.QueryRow(queryUserAccruals, userID)

	var accrualsUser int64
	if err := row.Scan(&accrualsUser); err != nil {
		// TODO Если у пользователя нет было, то Scan возвращает
		// error при записи NULL в int64 - юзать указатель
		l.logger.Err.Printf("could not get accruals: %s\n", err)
		return 0, nil
	}

	accrualsSum := float64(accrualsUser) / 100
	l.logger.Err.Printf("get accruals sum %f for user %s\n", accrualsSum, username)

	return accrualsSum, nil
}

func (l Loyalty) HowMatchUsed(username string) (float64, error) {

	var userID int64
	row := l.db.QueryRow(queryGetUserIDByName, username)
	if err := row.Scan(&userID); err != nil {
		return 0, market.ErrUserNotFound
	}

	return 0, nil
}

// SetAccrual - Изменение начислений по заказу
func (l Loyalty) SetAccrual(order int64, accrual float64) error {

	accrualRound := int64(accrual * 100)
	l.logger.Err.Printf("set accrual %d for order %d\n", accrualRound, order)
	_, err := l.db.Exec(queryUpdateAccrual, accrualRound, order)
	return err
}
