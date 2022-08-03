package service

import (
	"gophermarket/pkg/repository"
	"gophermarket/pkg/service/auth"
	"gophermarket/pkg/service/order"
)

type Service struct {
	Auth  auth.ServiceAuth
	Order order.ServiceOrder
}

func NewService(repo repository.Repository) *Service {

	return &Service{
		Auth:  auth.NewService(repo),
		Order: order.NewService(repo),
	}
}
