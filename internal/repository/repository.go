//go:generate mockgen -source repository.go -destination repository_mock.go -package repository
package repository

import (
	"context"

	market "gophermarket/internal"
)

const (
	StatusNew        = "NEW"
	StatusProcessing = "PROCESSING"
	StatusProcessed  = "PROCESSED"
	StatusInvalid    = "INVALID"
)

type Authorization interface {
	ID(ctx context.Context, user market.User) (int64, error)
	Create(ctx context.Context, user market.User) error
}

type Order interface {
	Create(ctx context.Context, number, username, status string) error
	CreateWithPayment(ctx context.Context, number, username string, sum float64) error
	SetStatus(ctx context.Context, order, status string) error
	UserOrders(ctx context.Context, username string) ([]OrderInfo, error)
	GetByStatuses(ctx context.Context, statuses []string) (map[string]string, error)
}

type Loyalty interface {
	SetAccrual(ctx context.Context, order string, accrual float64) error
	HowMatchUsed(ctx context.Context, username string) (float64, error)
	HowMatchAvailable(ctx context.Context, username string) (float64, error)
	Payments(ctx context.Context, username string) ([]PaymentInfo, error)
}

type Repository struct {
	Authorization
	Order
	Loyalty
}

// OrderLoyalty Информация о заказе из системы лояльности
type OrderLoyalty struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

// OrderInfo Информация о заказе
type OrderInfo struct {
	Order      string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

// PaymentInfo Информация о списании баллов
type PaymentInfo struct {
	OrderNum   string  `json:"order"`
	Sum        float64 `json:"sum"`
	UploadedAt string  `json:"processed_at"`
}

// OrderWithPay Заказ с оплатой баллами
type OrderWithPay struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
