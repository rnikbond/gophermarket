//go:generate mockgen -source repository.go -destination repository_mock.go -package repository
package repository

import (
	"context"

	market "gophermarket/internal"
	pkgOrder "gophermarket/pkg/order"
)

type Authorization interface {
	ID(ctx context.Context, user market.User) (int64, error)
	Create(ctx context.Context, user market.User) error
}

type Order interface {
	Create(ctx context.Context, number int64, username string, status string) error
	CreateWithPayment(ctx context.Context, number int64, username string, sum float64) error
	SetStatus(ctx context.Context, order int64, status string) error
	UserOrders(ctx context.Context, username string) ([]pkgOrder.InfoOrder, error)
	GetByStatuses(ctx context.Context, statuses []string) (map[int64]string, error)
}

type Loyalty interface {
	SetAccrual(ctx context.Context, order int64, accrual float64) error
	HowMatchUsed(ctx context.Context, username string) (float64, error)
	HowMatchAvailable(ctx context.Context, username string) (float64, error)
	WriteOffInfo(ctx context.Context, username string) ([]pkgOrder.WriteOff, error)
}

type Repository struct {
	Authorization
	Order
	Loyalty
}
