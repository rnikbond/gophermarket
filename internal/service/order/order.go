package order

import (
	"gophermarket/internal/repository"
	market "gophermarket/pkg"
)

type ServiceOrder interface {
	Create(number int64, username string) error
	Accruals(username string) (float64, error)
	Withdrawals(username string) (float64, error)
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

func (or Order) Accruals(username string) (float64, error) {

	return or.repo.Order.Accruals(username)
}

func (or Order) Withdrawals(username string) (float64, error) {

	return or.repo.Order.Withdrawals(username)
}
