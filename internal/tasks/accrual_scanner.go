package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"gophermarket/internal/repository"

	"github.com/sirupsen/logrus"
)

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusProcessed  = "PROCESSED"
	OrderStatusInvalid    = "INVALID"
)

type Accrual interface {
	Scan(ctx context.Context) error
}

type AccrualScanner struct {
	accrualAddr string
	repository  *repository.Repository
}

type Order struct {
	Order  string `json:"order"`
	Status string `json:"status"`
}

func NewScanner(addr string, repo *repository.Repository) Accrual {
	return &AccrualScanner{
		accrualAddr: addr,
		repository:  repo,
	}
}

func (scan AccrualScanner) Scan(ctx context.Context) error {

	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			order, errReload := scan.reloadOrders(ctx)
			if errReload != nil {
				logrus.Errorf("error reload orders from repository: %v\n", errReload)
				continue
			}

			scan.updateOrderStatuses(ctx, order)

		case <-ctx.Done():
			return nil
		}
	}
}

// updateOrderStatuses - Обновление статусов заказов в репозитории
func (scan AccrualScanner) updateOrderStatuses(ctx context.Context, orders []int64) {

	ordersData := scan.orderStatusesAccrual(ctx, orders)
	for _, order := range ordersData {
		if order.Status == OrderStatusProcessing {
			continue
		}

		if err := scan.setStatusOrder(ctx, order); err != nil {
			logrus.Errorf("error update status order in repository: %v\n", err)
		} else {
			fmt.Println("success update order status")
		}
	}
}

func (scan AccrualScanner) setStatusOrder(ctx context.Context, order Order) error {

	// TODO заюзать ctx

	orderNum, err := strconv.ParseInt(order.Order, 10, 64)
	if err != nil {
		return err
	}

	return scan.repository.Order.SetStatus(orderNum, order.Status)
}

// orderStatusesAccrual - Загрузка статусов по заказам из системы лояльности
func (scan AccrualScanner) orderStatusesAccrual(ctx context.Context, orders []int64) []Order {

	var ordersData []Order

	for _, orderNum := range orders {
		order, err := scan.orderStatusAccrual(ctx, orderNum)
		if err != nil {
			logrus.Printf("error load order from accrual: %v\n", err)
		}

		ordersData = append(ordersData, order)
	}

	return ordersData
}

// orderStatusAccrual - Загрузка статуса заказа из системы лояльности
func (scan AccrualScanner) orderStatusAccrual(ctx context.Context, orderNum int64) (Order, error) {

	// TODO заюзать ctx

	resp, errRequest := http.Get("http://" + scan.accrualAddr + "/api/orders/" + strconv.FormatInt(orderNum, 10))
	if errRequest != nil {
		return Order{}, errRequest
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.Errorf(err.Error())
		}
	}()

	data, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return Order{}, errRead
	}

	var order Order
	if err := json.Unmarshal(data, &order); err != nil {
		return Order{}, err
	}

	if len(order.Status) < 1 {
		return Order{}, errors.New("accrual service returned empty status")
	}

	return order, nil
}

// reloadOrders - Загрузка заказов из репозитория, у которых незавершенный статус начисления баллов
func (scan AccrualScanner) reloadOrders(ctx context.Context) ([]int64, error) {

	// TODO заюзать ctx

	orders, err := scan.repository.Order.GetByStatus(OrderStatusProcessing)
	if err != nil {
		logrus.Printf("repo return error : %v\n", err)
		return nil, errors.New(fmt.Sprintf("error reload processing orders: %v\n", err))
	}

	return orders, nil
}
