package order

import (
	market "gophermarket/pkg"
	"gophermarket/pkg/repository"
)

type IOrder interface {
	Create(number int64) error
}

type Order struct {
	repo repository.IRepository
}

func NewService(repo repository.IRepository) *Order {
	return &Order{
		repo: repo,
	}
}

func (or Order) Create(number int64) error {

	if !ValidOrder(number) {
		return market.ErrInvalidOrderNumber
	}

	return nil
}
