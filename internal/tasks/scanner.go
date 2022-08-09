package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"gophermarket/internal/repository"
	"gophermarket/pkg/logpack"
	pkgOrder "gophermarket/pkg/order"
)

type LoyaltyTask interface {
	Scan(ctx context.Context) error
}

type LoyaltyScanner struct {
	addr       string
	logger     *logpack.LogPack
	repository *repository.Repository
	client     *http.Client
}

func NewScanner(addr string, repo *repository.Repository, logger *logpack.LogPack) LoyaltyTask {
	return &LoyaltyScanner{
		addr:       addr,
		repository: repo,
		logger:     logger,
		client:     http.DefaultClient,
	}
}

func (scan LoyaltyScanner) Scan(ctx context.Context) error {

	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			orders, errReload := scan.reloadOrders(ctx)
			if errReload != nil {
				scan.logger.Err.Printf("error reload orders from repository: %s\n", errReload)
				continue
			}

			scan.updateOrderStatuses(ctx, orders)

		case <-ctx.Done():
			return nil
		}
	}
}

// updateOrderStatuses - Обновление статусов заказов в репозитории
func (scan LoyaltyScanner) updateOrderStatuses(ctx context.Context, orders map[int64]string) {

	ordersAccrual := scan.ordersAccrualService(ctx, orders)

	searcher := func(orders *[]pkgOrder.AccrualOrder, targetOrder int64) (int, bool) {
		for i, order := range *orders {
			num, err := strconv.ParseInt(order.Order, 10, 64)
			if err != nil {
				continue
			}

			if num == targetOrder {
				return i, true
			}
		}

		return 0, false
	}

	for orderNum, status := range orders {

		idx, ok := searcher(&ordersAccrual, orderNum)

		if !ok { // Незавершенный заказ, который есть в репозитории, не найден в системе лояльности
			if err := scan.repository.Order.SetStatus(ctx, orderNum, pkgOrder.StatusInvalid); err != nil {
				scan.logger.Err.Printf("error update status order in repository on %s: %s\n", pkgOrder.StatusInvalid, err)
			}
			continue
		}

		orderAccrual := ordersAccrual[idx]

		// Статус в репозитории такой же, как и в системе лояльности
		if status == orderAccrual.Status {
			continue
		}

		if err := scan.repository.Order.SetStatus(ctx, orderNum, orderAccrual.Status); err != nil {
			scan.logger.Err.Printf("error update status order in repository: %s\n", err)
		}

		if orderAccrual.Status != pkgOrder.StatusProcessed {
			continue
		}

		// Заказ получил завершенный статус в системе лояльности - сохраняем баллы за заказ
		if err := scan.repository.Loyalty.SetAccrual(ctx, orderNum, orderAccrual.Accrual); err != nil {
			scan.logger.Err.Printf("error update status order in repository: %s\n", err)
		}
	}
}

// orderStatusesAccrual - Загрузка статусов по заказам из системы лояльности
func (scan LoyaltyScanner) ordersAccrualService(ctx context.Context, orders map[int64]string) []pkgOrder.AccrualOrder {

	var ordersAccrual []pkgOrder.AccrualOrder

	for orderNum := range orders {
		orderAccrual, err := scan.orderAccrualService(ctx, orderNum)
		if err != nil {
			scan.logger.Err.Printf("error load order from accrual: %s\n", err)
			continue
		}

		ordersAccrual = append(ordersAccrual, orderAccrual)
	}

	return ordersAccrual
}

// orderAccrualService - Загрузка статуса заказа из системы лояльности
func (scan LoyaltyScanner) orderAccrualService(ctx context.Context, orderNum int64) (pkgOrder.AccrualOrder, error) {

	url := scan.addr + "/api/orders/" + strconv.FormatInt(orderNum, 10)
	request, errRequest := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if errRequest != nil {
		return pkgOrder.AccrualOrder{}, errRequest
	}

	resp, err := scan.client.Do(request)
	if err != nil {
		return pkgOrder.AccrualOrder{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return pkgOrder.AccrualOrder{}, errors.New("order not found")
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			scan.logger.Err.Printf("could not close response body: %s\n", err)
		}
	}()

	data, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return pkgOrder.AccrualOrder{}, errRead
	}

	var orderAccrual pkgOrder.AccrualOrder
	if err := json.Unmarshal(data, &orderAccrual); err != nil {
		return pkgOrder.AccrualOrder{}, err
	}

	if len(orderAccrual.Status) < 1 {
		return pkgOrder.AccrualOrder{}, errors.New("accrual service returned empty status")
	}

	return orderAccrual, nil
}

// reloadOrders - Загрузка заказов из репозитория, у которых незавершенный статус начисления баллов
func (scan LoyaltyScanner) reloadOrders(ctx context.Context) (map[int64]string, error) {

	statuses := []string{
		pkgOrder.StatusNew,
		pkgOrder.StatusProcessing,
	}
	orders, err := scan.repository.Order.GetByStatuses(ctx, statuses)
	if err != nil {
		scan.logger.Err.Printf("could not reload orders with statuses [%s,%s]: %s\n",
			pkgOrder.StatusNew, pkgOrder.StatusProcessing, err)
		return nil, err
	}

	return orders, nil
}
