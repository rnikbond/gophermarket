package service

import (
	"gophermarket/pkg/repository"
	"gophermarket/pkg/service/auth"
	"gophermarket/pkg/service/order"
)

type Service struct {
	auth.IAuthorization
	order.IOrder
}

func NewService(repo repository.IRepository) *Service {

	return &Service{
		IAuthorization: auth.NewService(repo),
		IOrder:         order.NewService(repo),
	}
}
