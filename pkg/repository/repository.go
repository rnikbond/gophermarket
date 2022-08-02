//go:generate mockgen -source repository.go -destination ../../mocks/pkg/repository/repository.go -package mock_repository
package repository

import (
	market "gophermarket/pkg"
)

type IRepository interface {
	SignUp(user market.User) error
	SignIn(user market.User) error

	CreateOrder(number int64) error

	Finish() error
}
