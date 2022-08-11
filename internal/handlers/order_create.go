package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	market "gophermarket/pkg"
)

// CreateOrder Обработчик запроса на создание заказа
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	username, _, ok := r.BasicAuth()
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			h.logger.Err.Printf("could not close request body: %s\n", err)
		}
	}()

	order, err := io.ReadAll(r.Body)
	if len(order) == 0 || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	errOrder := h.services.Order.Create(r.Context(), string(order), username)
	if errOrder == nil {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	if errOrder != market.ErrUserAlreadyOrderedIt {
		http.Error(w, errOrder.Error(), market.ErrorHTTP(errOrder))
		return
	}
}

// CreateWithPay - обработчик запроса на создание заказа со списанием баллов
func (h *Handler) CreateWithPay(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	username, _, ok := r.BasicAuth()
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			h.logger.Err.Printf("could not close request body: %s\n", err)
		}
	}()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var orderPay market.OrderWithPay

	if err := json.Unmarshal(data, &orderPay); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	errCreate := h.services.Order.CreateWithPayment(r.Context(), orderPay.Order, username, orderPay.Sum)
	if errCreate != nil {
		http.Error(w, errCreate.Error(), market.ErrorHTTP(errCreate))
		return
	}
}
