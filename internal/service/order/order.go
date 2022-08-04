package order

import (
	"gophermarket/internal/repository"
	market "gophermarket/pkg"
)

type ServiceOrder interface {
	Create(number int64, username string) error
}

type Order struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) ServiceOrder {
	return &Order{
		repo: repo,
	}
}

func (or Order) Create(number int64, username string) error {

	if !ValidOrder(number) {
		return market.ErrInvalidOrderNumber
	}

	return or.repo.Order.Create(number, username)
}
