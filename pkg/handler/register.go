package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	market "gophermarket/pkg"

	"github.com/sirupsen/logrus"
)

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			logrus.Errorf(fmt.Sprintf("error close body request: %v\n", err))
		}
	}()

	var user market.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if h.services.Registration == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.services.Registration.CreateUser(user); err != nil {
		http.Error(w, err.Error(), market.ErrorHTTP(err))
		return
	}
}
