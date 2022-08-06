//go:generate mockgen -source repository.go -destination repository_mock.go -package repository
package repository

import (
	market "gophermarket/internal"
)

type Authorization interface {
	Create(user market.User) error
	ID(user market.User) (int64, error)
}

type Order interface {
	Create(number int64, username string) error
	GetByStatus(status string) ([]int64, error)
	SetStatus(order int64, status string) error
}

type Repository struct {
	Authorization
	Order
}
