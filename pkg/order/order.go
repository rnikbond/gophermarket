package order

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

type AccrualOrder struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

type Order struct {
	Order  int64
	Status string
}

type InfoOrder struct {
	Order      string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

func (o AccrualOrder) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Order  : %s\n", o.Order))
	builder.WriteString(fmt.Sprintf("Status : %s\n", o.Status))
	builder.WriteString(fmt.Sprintf("Accrual: %f\n", o.Accrual))

	return builder.String()
}
