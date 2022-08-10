//go:generate mockgen -source variant.go -destination order_mock.go -package order
package order

import (
	"context"
	"strconv"

	"gophermarket/internal/repository"
	"gophermarket/pkg"
	"gophermarket/pkg/logpack"

	"github.com/EClaesson/go-luhn"
)

type ServiceOrder interface {
	Create(ctx context.Context, number int64, username string) error
	CreateWithPayment(ctx context.Context, number int64, username string, sum float64) error
	UserOrders(ctx context.Context, username string) ([]pkg.OrderInfo, error)
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

	if ok, err := luhn.IsValid(strconv.FormatInt(number, 10)); !ok || err != nil {
		if err != nil {
			or.logger.Err.Printf("could not validate order number: %s\n", err)
		}
		return pkg.ErrInvalidOrderNumber
	}

	return or.repo.Order.Create(ctx, number, username, pkg.StatusNew)
}

func (or Order) CreateWithPayment(ctx context.Context, number int64, username string, sum float64) error {

	if ok, err := luhn.IsValid(strconv.FormatInt(number, 10)); !ok || err != nil {
		if err != nil {
			or.logger.Err.Printf("could not validate order number: %s\n", err)
		}
		return pkg.ErrInvalidOrderNumber
	}

	current, errCurrent := or.repo.Loyalty.HowMatchAvailable(ctx, username)
	if errCurrent != nil {
		return errCurrent
	}

	used, errUsed := or.repo.Loyalty.HowMatchUsed(ctx, username)
	if errUsed != nil {
		return errCurrent
	}

	if (current - used) < sum {
		return pkg.ErrPaymentNotAvailable
	}

	return or.repo.Order.CreateWithPayment(ctx, number, username, sum)
}

func (or Order) UserOrders(ctx context.Context, username string) ([]pkg.OrderInfo, error) {
	return or.repo.Order.UserOrders(ctx, username)
}
