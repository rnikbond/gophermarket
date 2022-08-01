package service

import (
	market "gophermarket/pkg"
	"gophermarket/pkg/repository"
)

type Authorization interface {
	SignUp(user market.User) (string, error)
	SignIn(user market.User) (string, error)
}

type Service struct {
	Authorization
}

func NewService(repo *repository.Repository) *Service {

	return &Service{
		Authorization: NewAuthService(repo),
	}
}
