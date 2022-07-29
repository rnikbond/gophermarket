package service

import (
	"gophermarket/pkg/repository"
)

type Authorization interface {
}

type Registration interface {
}

type Service struct {
	Authorization
	Registration
}

func NewService(repos *repository.Repository) *Service {

	return &Service{}
}
