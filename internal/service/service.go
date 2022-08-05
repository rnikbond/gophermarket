package service

import (
	"gophermarket/internal/repository"
	"gophermarket/internal/service/auth"
	"gophermarket/internal/service/order"
)

type Service struct {
	Auth  auth.ServiceAuth
	Order order.ServiceOrder
}

func NewServices(repo *repository.Repository, pwdSalt string) *Service {

	return &Service{
		Auth:  auth.NewService(repo, pwdSalt),
		Order: order.NewService(repo),
	}
}
