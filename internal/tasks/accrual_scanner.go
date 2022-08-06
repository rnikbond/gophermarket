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
	pkgOrder "gophermarket/pkg/order"

	"github.com/sirupsen/logrus"
)

type Accrual interface {
	Scan(ctx context.Context) error
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

func (scan AccrualScanner) Scan(ctx context.Context) error {

	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			orders, errReload := scan.reloadOrders(ctx)
			if errReload != nil {
				logrus.Errorf("error reload orders from repository: %v\n", errReload)
				continue
			}

			scan.updateOrderStatuses(ctx, orders)

		case <-ctx.Done():
			return nil
		}
	}
}

// updateOrderStatuses - Обновление статусов заказов в репозитории
func (scan AccrualScanner) updateOrderStatuses(ctx context.Context, orders map[int64]string) {

	ordersAccrual := scan.orderStatusesAccrual(ctx, orders)

	for orderNum, status := range orders {

		statusAccrual, ok := ordersAccrual[orderNum]

		if !ok { // Незавершенный заказ, который есть в репозитории, не найден в системе лояльности

			if err := scan.repository.Order.SetStatus(orderNum, pkgOrder.StatusInvalid); err != nil {
				logrus.Errorf("error update status order in repository on %s: %v\n", pkgOrder.StatusInvalid, err)
			}
			continue
		}

		if status == statusAccrual {
			continue
		}

		if err := scan.repository.Order.SetStatus(orderNum, statusAccrual); err != nil {
			logrus.Errorf("error update status order in repository: %v\n", err)
		}
	}
}

// orderStatusesAccrual - Загрузка статусов по заказам из системы лояльности
func (scan AccrualScanner) orderStatusesAccrual(ctx context.Context, orders map[int64]string) map[int64]string {

	ordersAccrual := make(map[int64]string)

	for orderNum := range orders {
		statusAccrual, err := scan.orderStatusAccrual(ctx, orderNum)
		if err != nil {
			logrus.Printf("error load order from accrual: %v\n", err)
			continue
		}

		ordersAccrual[orderNum] = statusAccrual
	}

	return ordersAccrual
}

// orderStatusAccrual - Загрузка статуса заказа из системы лояльности
func (scan AccrualScanner) orderStatusAccrual(ctx context.Context, orderNum int64) (string, error) {

	// TODO заюзать ctx

	resp, errRequest := http.Get("http://" + scan.accrualAddr + "/api/orders/" + strconv.FormatInt(orderNum, 10))
	if errRequest != nil {
		return ``, errRequest
	}

	if resp.StatusCode != http.StatusOK {
		return ``, errors.New("order not found")
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logrus.Errorf(err.Error())
		}
	}()

	data, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return ``, errRead
	}

	var orderAccrual pkgOrder.AccrualOrder
	if err := json.Unmarshal(data, &orderAccrual); err != nil {
		return ``, err
	}

	if len(orderAccrual.Status) < 1 {
		return ``, errors.New("accrual service returned empty status")
	}

	return orderAccrual.Status, nil
}

// reloadOrders - Загрузка заказов из репозитория, у которых незавершенный статус начисления баллов
func (scan AccrualScanner) reloadOrders(ctx context.Context) (map[int64]string, error) {

	// TODO заюзать ctx

	statuses := []string{
		pkgOrder.StatusNew,
		pkgOrder.StatusProcessing,
	}
	orders, err := scan.repository.Order.GetByStatuses(statuses)
	if err != nil {
		logrus.Printf("repo return error : %v\n", err)
		return nil, errors.New(fmt.Sprintf("error reload processing orders: %v\n", err))
	}

	return orders, nil
}
