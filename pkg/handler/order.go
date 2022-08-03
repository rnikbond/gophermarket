package handler

import (
	"io"
	"net/http"
	"strconv"

	market "gophermarket/pkg"
	"gophermarket/pkg/service/auth"

	"github.com/sirupsen/logrus"
)

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	username, err := auth.ParseCookie(r)
	if err != nil {
		http.Error(w, err.Error(), market.ErrorHTTP(err))
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

	number, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		http.Error(w, "invalid order number", http.StatusBadRequest)
		return
	}

	if err := h.services.Order.Create(number, username); err != nil {
		http.Error(w, err.Error(), market.ErrorHTTP(err))
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
