package handler

import (
	"io"
	"net/http"
	"strconv"

	market "gophermarket/pkg"

	"github.com/sirupsen/logrus"
)

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {

	username, _, ok := r.BasicAuth()
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			logrus.Errorf("handler CreateOrder :: error close request body: %v", err)
		}
	}()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		http.Error(w, "invalid order number", http.StatusBadRequest)
		return
	}

	errOrder := h.services.Order.Create(order, username)
	if errOrder == nil {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	if errOrder != market.ErrUserAlreadyOrderedIt {
		http.Error(w, errOrder.Error(), market.ErrorHTTP(errOrder))
		return
	}
}
