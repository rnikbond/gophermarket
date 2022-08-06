package order

import (
	"strconv"
)

const (
	StatusNew        = "NEW"
	StatusRegistered = "REGISTERED"
	StatusProcessing = "PROCESSING"
	StatusProcessed  = "PROCESSED"
	StatusInvalid    = "INVALID"
)

type AccrualOrder struct {
	Order  string `json:"order"`
	Status string `json:"status"`
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
