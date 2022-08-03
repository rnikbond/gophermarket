//go:generate mockgen -source repository.go -destination ../../mocks/pkg/repository/repository.go -package mock_repository
package repository

import (
	market "gophermarket/pkg"
)

type Repository interface {
	SignUp(user market.User) error
	SignIn(user market.User) error

	CreateOrder(number int64, username string) error

	Finish() error
}
