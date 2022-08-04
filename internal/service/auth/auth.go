package auth

import (
	market "gophermarket/internal"
	"gophermarket/internal/repository"
	"gophermarket/pkg"
)

type ServiceAuth interface {
	SignUp(user market.User) error
	SignIn(user market.User) error

	ValidateAuth(user market.User) error
}

type Auth struct {
	repo         *repository.Repository
	passwordSalt string
}

func NewService(repo *repository.Repository, pwdSalt string) ServiceAuth {
	return &Auth{
		repo:         repo,
		passwordSalt: pwdSalt,
	}
}

func (s Auth) SignUp(user market.User) error {
	if err := s.ValidateAuth(user); err != nil {
		return err
	}

	hash, err := s.GeneratePasswordHash(user.Password)
	if err != nil {
		return err
	}

	user.Password = hash

	if err := s.repo.Authorization.Create(user); err != nil {
		return err
	}

	return nil
}

func (s Auth) SignIn(user market.User) error {
	if err := s.ValidateAuth(user); err != nil {
		return err
	}

	hash, err := s.GeneratePasswordHash(user.Password)
	if err != nil {
		return err
	}

	user.Password = hash

	if _, err := s.repo.Authorization.ID(user); err != nil {
		return err
	}

	return nil
}

func (s Auth) ValidateAuth(user market.User) error {

	if len(user.Username) < 1 || len(user.Password) < 1 {
		return pkg.ErrEmptyAuthData
	}

	return nil
}
