//go:generate mockgen -source loyalty.go -destination loyalty_mock.go -package loyalty
package loyalty

import (
	"context"
	"math"

	"gophermarket/internal/repository"
	"gophermarket/pkg/logpack"
)

type ServiceLoyalty interface {
	HowMatchAvailable(ctx context.Context, username string) (float64, error)
	HowMatchUsed(ctx context.Context, username string) (float64, error)
	SetAccrual(ctx context.Context, order string, accrual float64) error
	Balance(ctx context.Context, username string) (Balance, error)
	Payments(ctx context.Context, username string) ([]repository.PaymentInfo, error)
}

type Loyalty struct {
	logger *logpack.LogPack
	repo   *repository.Repository
}

type Balance struct {
	Accrual   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
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

func (service Loyalty) SetAccrual(ctx context.Context, order string, accrual float64) error {

	return service.repo.Loyalty.SetAccrual(ctx, order, accrual)
}

func (service Loyalty) Balance(ctx context.Context, username string) (Balance, error) {

	current, err := service.HowMatchAvailable(ctx, username)
	if err != nil {
		return Balance{}, err
	}

	used, err := service.HowMatchUsed(ctx, username)
	if err != nil {
		return Balance{}, err
	}

	round := func(val float64, precision uint) float64 {
		ratio := math.Pow(10, float64(precision))
		return math.Round(val*ratio) / ratio
	}

	current = round(current-used, 2)

	return Balance{
		Accrual:   current,
		Withdrawn: used,
	}, nil
}

func (service Loyalty) Payments(ctx context.Context, username string) ([]repository.PaymentInfo, error) {
	return service.repo.Loyalty.Payments(ctx, username)
}
