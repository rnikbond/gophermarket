package service

import (
	"crypto/sha256"
	"encoding/hex"

	market "gophermarket/pkg"
	"gophermarket/pkg/repository"
)

const salt = "ololo-ololo-ololo"

type RegisterService struct {
	repo repository.Registration
}

func NewRegisterService(repo repository.Registration) *RegisterService {
	return &RegisterService{repo: repo}
}

func (s *RegisterService) CreateUser(user market.User) error {

	if len(user.Username) < 1 || len(user.Password) < 1 {
		return market.ErrEmptyLoginPassword
	}

	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *RegisterService) generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum([]byte(salt)))
}
