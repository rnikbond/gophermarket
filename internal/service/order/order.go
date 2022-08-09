package order

import (
	"context"

	"gophermarket/internal/repository"
	market "gophermarket/pkg"
	"gophermarket/pkg/logpack"
	"gophermarket/pkg/order"
)

type ServiceOrder interface {
	Create(ctx context.Context, number int64, username string) error
	UserOrders(ctx context.Context, username string) ([]order.InfoOrder, error)
}

type Order struct {
	logger *logpack.LogPack
	repo   *repository.Repository
}

func NewService(repo *repository.Repository, logger *logpack.LogPack) ServiceOrder {
	return &Order{
		repo:   repo,
		logger: logger,
	}
}

func (or Order) Create(ctx context.Context, number int64, username string) error {

	if !ValidOrder(number) {
		return market.ErrInvalidOrderNumber
	}

	return or.repo.Order.Create(ctx, number, username)
}

func (or Order) UserOrders(ctx context.Context, username string) ([]order.InfoOrder, error) {
	return or.repo.Order.UserOrders(ctx, username)
}
