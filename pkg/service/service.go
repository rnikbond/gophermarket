package service

import (
	"gophermarket/pkg/repository"
	"gophermarket/pkg/service/auth"
)

type Service struct {
	auth.Authorization
}

func NewService(repo repository.Repository) *Service {

	return &Service{
		Authorization: auth.NewService(repo),
	}
}
