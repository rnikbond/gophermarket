package repository

import (
	gophermarket "gophermarket/pkg"

	"github.com/jmoiron/sqlx"
)

type LoginService struct {
	db *sqlx.DB
}

func NewLoginService(db *sqlx.DB) *LoginService {
	return &LoginService{db: db}
}

func (s *LoginService) Login(user gophermarket.User) error {

	return nil
}
