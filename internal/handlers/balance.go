package handlers

import (
	"encoding/json"
	"net/http"

	"gophermarket/pkg"
)

func (h *Handler) Balance(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	username, _, ok := r.BasicAuth()
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	balance, err := h.services.Loyalty.Balance(username)
	if err != nil {
		w.WriteHeader(pkg.ErrorHTTP(err))
		return
	}

	data, err := json.Marshal(&balance)
	if err != nil {
		http.Error(w, "error marshal balance", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(data); err != nil {
		http.Error(w, "error write json-balance in response", http.StatusInternalServerError)
		return
	}
}
