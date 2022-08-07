package order

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	StatusNew        = "NEW"
	StatusRegistered = "REGISTERED"
	StatusProcessing = "PROCESSING"
	StatusProcessed  = "PROCESSED"
	StatusInvalid    = "INVALID"
)

type AccrualOrder struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int64  `json:"accrual,omitempty"`
}

type Order struct {
	Order  int64
	Status string
}

func ToOrder(order AccrualOrder) (Order, error) {

	num, err := strconv.ParseInt(order.Order, 10, 64)
	if err != nil {
		return Order{}, err
	}

	return Order{
		Order:  num,
		Status: order.Status,
	}, nil
}

func (o AccrualOrder) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Order  : %s\n", o.Order))
	builder.WriteString(fmt.Sprintf("Status : %s\n", o.Status))
	builder.WriteString(fmt.Sprintf("Accrual: %d\n", o.Accrual))

	return builder.String()
}
