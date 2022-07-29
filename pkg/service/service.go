package service

import (
	market "gophermarket/pkg"
	"gophermarket/pkg/repository"
)

type Authorization interface {
	Login(user market.User) error
}

type Registration interface {
	CreateUser(user market.User) error
}

type Service struct {
	Authorization
	Registration
}

func NewService(repo *repository.Repository) *Service {

	return &Service{
		Registration: NewRegisterService(repo),
	}
}
