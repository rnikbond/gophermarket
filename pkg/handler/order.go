package handler

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	market "gophermarket/pkg"

	"github.com/sirupsen/logrus"
)

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			logrus.Errorf("handler CreateOrder :: error close request body: %v", err)
		}
	}()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	number, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.services.IOrder.Create(number); err != nil {
		http.Error(w, err.Error(), market.ErrorHTTP(err))
		return
	}

	fmt.Println(number)
}
