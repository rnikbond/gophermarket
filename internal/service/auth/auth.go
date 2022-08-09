package auth

import (
	"context"

	market "gophermarket/internal"
	"gophermarket/internal/repository"
	"gophermarket/pkg"
	"gophermarket/pkg/logpack"
)

type ServiceAuth interface {
	SignUp(ctx context.Context, user market.User) error
	SignIn(ctx context.Context, user market.User) error

	ValidateAuth(user market.User) error
}

type Auth struct {
	logger       *logpack.LogPack
	repo         *repository.Repository
	passwordSalt string
}

func NewService(repo *repository.Repository, pwdSalt string, logger *logpack.LogPack) ServiceAuth {
	return &Auth{
		repo:         repo,
		passwordSalt: pwdSalt,
		logger:       logger,
	}
}

func (s Auth) SignUp(ctx context.Context, user market.User) error {
	if err := s.ValidateAuth(user); err != nil {
		return err
	}

	hash, err := pkg.GeneratePasswordHash(user.Password, s.passwordSalt)
	if err != nil {
		return err
	}

	user.Password = hash

	if err := s.repo.Authorization.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s Auth) SignIn(ctx context.Context, user market.User) error {
	if err := s.ValidateAuth(user); err != nil {
		return err
	}

	hash, err := pkg.GeneratePasswordHash(user.Password, s.passwordSalt)
	if err != nil {
		return err
	}

	user.Password = hash

	if _, err := s.repo.Authorization.ID(ctx, user); err != nil {
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
