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

func (o AccrualOrder) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Order  : %s\n", o.Order))
	builder.WriteString(fmt.Sprintf("Status : %s\n", o.Status))
	builder.WriteString(fmt.Sprintf("Accrual: %d\n", o.Accrual))

	return builder.String()
}
