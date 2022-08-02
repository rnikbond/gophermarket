package auth

import (
	market "gophermarket/pkg"
	"gophermarket/pkg/repository"
)

type Authorization interface {
	SignUp(user market.User) (string, error)
	SignIn(user market.User) (string, error)

	ValidateAuth(user market.User) error
}

type Auth struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Auth {
	return &Auth{repo: repo}
}

func (s Auth) SignUp(user market.User) (string, error) {
	if err := s.ValidateAuth(user); err != nil {
		return ``, err
	}

	hash, err := GeneratePasswordHash(user.Password)
	if err != nil {
		return ``, err
	}

	user.Password = hash

	if err := s.repo.SignUp(user); err != nil {
		return ``, err
	}

	return GenerateJWT(user), nil
}

func (s Auth) SignIn(user market.User) (string, error) {
	if err := s.ValidateAuth(user); err != nil {
		return ``, err
	}

	hash, err := GeneratePasswordHash(user.Password)
	if err != nil {
		return ``, err
	}

	user.Password = hash

	if err := s.repo.SignIn(user); err != nil {
		return ``, err
	}

	return GenerateJWT(user), nil
}

func (s Auth) ValidateAuth(user market.User) error {

	if len(user.Username) < 1 || len(user.Password) < 1 {
		return market.ErrEmptyAuthData
	}

	return nil
}
