//go:generate mockgen -source loyalty.go -destination loyalty_mock.go -package loyalty
package loyalty

import (
	"context"

	gophermarket "gophermarket/internal"
	"gophermarket/internal/repository"
	"gophermarket/pkg/logpack"
)

type ServiceLoyalty interface {
	HowMatchAvailable(ctx context.Context, username string) (float64, error)
	HowMatchUsed(ctx context.Context, username string) (float64, error)
	SetAccrual(ctx context.Context, order int64, accrual float64) error
	Balance(ctx context.Context, username string) (gophermarket.Balance, error)
}

type Loyalty struct {
	logger *logpack.LogPack
	repo   *repository.Repository
}

func NewService(repo *repository.Repository, logger *logpack.LogPack) ServiceLoyalty {
	return &Loyalty{
		repo:   repo,
		logger: logger,
	}
}

func (service Loyalty) HowMatchAvailable(ctx context.Context, username string) (float64, error) {

	return service.repo.Loyalty.HowMatchAvailable(ctx, username)
}

func (service Loyalty) HowMatchUsed(ctx context.Context, username string) (float64, error) {

	return service.repo.Loyalty.HowMatchUsed(ctx, username)
}

func (service Loyalty) SetAccrual(ctx context.Context, order int64, accrual float64) error {

	return service.repo.Loyalty.SetAccrual(ctx, order, accrual)
}

func (service Loyalty) Balance(ctx context.Context, username string) (gophermarket.Balance, error) {

	current, err := service.HowMatchAvailable(ctx, username)
	if err != nil {
		return gophermarket.Balance{}, err
	}

	used, err := service.HowMatchUsed(ctx, username)
	if err != nil {
		return gophermarket.Balance{}, err
	}

	return gophermarket.Balance{
		Accrual:   current,
		Withdrawn: used,
	}, nil
}
