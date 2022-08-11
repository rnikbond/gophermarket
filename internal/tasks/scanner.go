package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"gophermarket/internal/repository"
	"gophermarket/pkg"
	"gophermarket/pkg/logpack"
)

type LoyaltyTask interface {
	Scan(ctx context.Context)
}

type LoyaltyScanner struct {
	addr       string
	logger     *logpack.LogPack
	repository *repository.Repository
	client     *http.Client
	interval   time.Duration
}

func NewScanner(addr string, repo *repository.Repository, interval time.Duration, logger *logpack.LogPack) LoyaltyTask {
	return &LoyaltyScanner{
		addr:       addr,
		repository: repo,
		logger:     logger,
		client:     http.DefaultClient,
		interval:   interval,
	}
}

func (scan LoyaltyScanner) Scan(ctx context.Context) {

	ticker := time.NewTicker(scan.interval)

	go func() {
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
				return
			}
		}
	}()
}

// updateOrderStatuses - Обновление статусов заказов в репозитории
func (scan LoyaltyScanner) updateOrderStatuses(ctx context.Context, orders map[string]string) {

	ordersAccrual := scan.ordersAccrualService(ctx, orders)

	searcher := func(orders *[]pkg.OrderLoyalty, targetOrder string) (int, bool) {
		for i, order := range *orders {
			if order.Order == targetOrder {
				return i, true
			}
		}

		return 0, false
	}

	for order, status := range orders {

		idx, ok := searcher(&ordersAccrual, order)

		if !ok { // Незавершенный заказ, который есть в репозитории, не найден в системе лояльности
			if err := scan.repository.Order.SetStatus(ctx, order, pkg.StatusInvalid); err != nil {
				scan.logger.Err.Printf("error update status order in repository on %s: %s\n", pkg.StatusInvalid, err)
			}
			continue
		}

		orderAccrual := ordersAccrual[idx]

		// Статус в репозитории такой же, как и в системе лояльности
		if status == orderAccrual.Status {
			continue
		}

		if err := scan.repository.Order.SetStatus(ctx, order, orderAccrual.Status); err != nil {
			scan.logger.Err.Printf("error update status order in repository: %s\n", err)
		}

		if orderAccrual.Status != pkg.StatusProcessed {
			continue
		}

		// Заказ получил завершенный статус в системе лояльности - сохраняем баллы за заказ
		if err := scan.repository.Loyalty.SetAccrual(ctx, order, orderAccrual.Accrual); err != nil {
			scan.logger.Err.Printf("error update status order in repository: %s\n", err)
		}
	}
}

// orderStatusesAccrual - Загрузка статусов по заказам из системы лояльности
func (scan LoyaltyScanner) ordersAccrualService(ctx context.Context, orders map[string]string) []pkg.OrderLoyalty {

	var ordersAccrual []pkg.OrderLoyalty

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
func (scan LoyaltyScanner) orderAccrualService(ctx context.Context, order string) (pkg.OrderLoyalty, error) {

	url := scan.addr + "/api/orders/" + order
	request, errRequest := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if errRequest != nil {
		return pkg.OrderLoyalty{}, errRequest
	}

	resp, err := scan.client.Do(request)
	if err != nil {
		return pkg.OrderLoyalty{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return pkg.OrderLoyalty{}, errors.New("order not found")
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			scan.logger.Err.Printf("could not close response body: %s\n", err)
		}
	}()

	data, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return pkg.OrderLoyalty{}, errRead
	}

	var orderAccrual pkg.OrderLoyalty
	if err := json.Unmarshal(data, &orderAccrual); err != nil {
		return pkg.OrderLoyalty{}, err
	}

	if len(orderAccrual.Status) == 0 {
		return pkg.OrderLoyalty{}, errors.New("accrual service returned empty status")
	}

	return orderAccrual, nil
}

// reloadOrders - Загрузка заказов из репозитория, у которых незавершенный статус начисления баллов
func (scan LoyaltyScanner) reloadOrders(ctx context.Context) (map[string]string, error) {

	statuses := []string{
		pkg.StatusNew,
		pkg.StatusProcessing,
	}
	orders, err := scan.repository.Order.GetByStatuses(ctx, statuses)
	if err != nil {
		scan.logger.Err.Printf("could not reload orders with statuses [%s,%s]: %s\n",
			pkg.StatusNew, pkg.StatusProcessing, err)
		return nil, err
	}

	return orders, nil
}
