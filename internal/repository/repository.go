//go:generate mockgen -source repository.go -destination mock/repository.go -package mock_repository
package repository

import (
	market "gophermarket/internal"
)

type Authorization interface {
	Create(user market.User) error
	Id(user market.User) (int64, error)
}

type Order interface {
	Create(number int64, username string) error
}

type Repository struct {
	Authorization
	Order
}
