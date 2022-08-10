package service

import (
	"gophermarket/internal/repository"
	"gophermarket/internal/service/auth"
	"gophermarket/internal/service/loyalty"
	"gophermarket/internal/service/order"
	"gophermarket/pkg/logpack"
)

type Service struct {
	Auth    auth.ServiceAuth
	Order   order.ServiceOrder
	Loyalty loyalty.ServiceLoyalty
}

func NewServices(repo *repository.Repository, pwdSalt string, logger *logpack.LogPack) *Service {

	return &Service{
		Auth:    auth.NewService(repo, pwdSalt, logger),
		Order:   order.NewService(repo, logger),
		Loyalty: loyalty.NewService(repo, logger),
	}
}
