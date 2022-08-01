package service

import (
	market "gophermarket/pkg"
	"gophermarket/pkg/repository"
)

const pwdSalt = "ololo-ololo-ololo"

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) SignIn(user market.User) (string, error) {
	if len(user.Username) < 1 || len(user.Password) < 1 {
		return "", market.ErrEmptyLoginPassword
	}

	user.Password = s.generatePasswordHash(user)

	if err := s.repo.SignIn(user); err != nil {
		return "", err
	}

	return s.generateJWT(user), nil
}

func (s *AuthService) SignUp(user market.User) (string, error) {

	if len(user.Username) < 1 || len(user.Password) < 1 {
		return "", market.ErrEmptyLoginPassword
	}

	user.Password = s.generatePasswordHash(user)

	if err := s.repo.SignUp(user); err != nil {
		return "", err
	}

	user.Token = s.generateJWT(user)

	return s.generateJWT(user), nil
}
