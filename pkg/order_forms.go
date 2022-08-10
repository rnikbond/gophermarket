package pkg

import (
	"fmt"
	"strings"
)

const (
	StatusNew        = "NEW"
	StatusProcessing = "PROCESSING"
	StatusProcessed  = "PROCESSED"
	StatusInvalid    = "INVALID"
)

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

func (o OrderLoyalty) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Order  : %s\n", o.Order))
	builder.WriteString(fmt.Sprintf("Status : %s\n", o.Status))
	builder.WriteString(fmt.Sprintf("Accrual: %f\n", o.Accrual))

	return builder.String()
}
