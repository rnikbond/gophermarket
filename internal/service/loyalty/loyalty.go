package loyalty

import (
	gophermarket "gophermarket/internal"
	"gophermarket/internal/repository"
	"gophermarket/pkg/logpack"
)

type ServiceLoyalty interface {
	HowMatchAvailable(username string) (float64, error)
	HowMatchUsed(username string) (float64, error)
	SetAccrual(order int64, accrual float64) error
	Balance(username string) (gophermarket.Balance, error)
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

func (service Loyalty) HowMatchAvailable(username string) (float64, error) {

	return service.repo.Loyalty.HowMatchAvailable(username)
}

func (service Loyalty) HowMatchUsed(username string) (float64, error) {

	return service.repo.Loyalty.HowMatchUsed(username)
}

func (service Loyalty) SetAccrual(order int64, accrual float64) error {

	return service.repo.Loyalty.SetAccrual(order, accrual)
}

func (service Loyalty) Balance(username string) (gophermarket.Balance, error) {

	current, err := service.HowMatchAvailable(username)
	if err != nil {
		return gophermarket.Balance{}, nil
	}

	used, err := service.HowMatchUsed(username)
	if err != nil {
		return gophermarket.Balance{}, nil
	}

	return gophermarket.Balance{
		Accrual:   current,
		Withdrawn: used,
	}, nil
}
