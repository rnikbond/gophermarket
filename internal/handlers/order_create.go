package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	market "gophermarket/pkg"
)

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

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	errOrder := h.services.Order.Create(r.Context(), order, username)
	if errOrder == nil {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	if errOrder != market.ErrUserAlreadyOrderedIt {
		http.Error(w, errOrder.Error(), market.ErrorHTTP(errOrder))
		return
	}
}

// CreateWithPay - ќбработчик запроса на создание заказа со списанием баллов
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

	orderPay := struct {
		Order string  `json:"order"`
		Sum   float64 `json:"sum"`
	}{}

	if err := json.Unmarshal(data, &orderPay); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	orderNum, errParse := strconv.ParseInt(orderPay.Order, 10, 64)
	if errParse != nil {
		http.Error(w, errParse.Error(), http.StatusBadRequest)
		return
	}

	errCreate := h.services.Order.CreateWithPayment(r.Context(), orderNum, username, orderPay.Sum)
	if errCreate != nil {
		http.Error(w, errCreate.Error(), market.ErrorHTTP(errCreate))
		return
	}
}
