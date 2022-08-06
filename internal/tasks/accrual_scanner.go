package tasks

import (
	"context"

	"gophermarket/internal/repository"
)

type Accrual interface {
	Scan(ctx *context.Context) error
}

type AccrualScanner struct {
	accrualAddr string
	repository  *repository.Repository
}

func NewScanner(addr string, repo *repository.Repository) Accrual {
	return &AccrualScanner{
		accrualAddr: addr,
		repository:  repo,
	}
}

func (scan AccrualScanner) Scan(ctx *context.Context) error {

	//for {
	//	select {
	//
	//	}
	//}

	return nil
}
