//go:generate mockgen -source repository.go -destination repository_mock.go -package repository
package repository

import (
	market "gophermarket/internal"
	pkgOrder "gophermarket/pkg/order"
)

type Authorization interface {
	Create(user market.User) error
	ID(user market.User) (int64, error)
}

type Order interface {
	Create(number int64, username string) error
	GetByStatuses(statuses []string) (map[int64]string, error)
	SetStatus(order int64, status string) error
	UserOrders(username string) ([]pkgOrder.InfoOrder, error)
}

type Loyalty interface {
	HowMatchAvailable(username string) (float64, error)
	HowMatchUsed(username string) (float64, error)
	SetAccrual(order int64, accrual float64) error
}

type Repository struct {
	Authorization
	Order
	Loyalty
}
