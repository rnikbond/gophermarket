package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"gophermarket/pkg"
)

func (h *Handler) OrdersList(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	username, _, ok := r.BasicAuth()
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	orders, err := h.services.Order.UserOrders(username)
	if err != nil {
		log.Printf("OrdersList :: error get user orders list: %v\n", err)
		w.WriteHeader(pkg.ErrorHTTP(err))
		return
	}

	if len(orders) < 1 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	data, err := json.Marshal(&orders)
	if err != nil {
		log.Printf("OrdersList :: error marshal user orders list: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(data); err != nil {
		log.Printf("OrdersList :: error write to response marshaled user orders list: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
