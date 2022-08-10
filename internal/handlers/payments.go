package handlers

import (
	"encoding/json"
	"net/http"

	"gophermarket/pkg"
)

// Payments - Обработчик запроса данных о списании баллов
func (h *Handler) Payments(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	username, _, ok := r.BasicAuth()
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	info, err := h.services.Loyalty.Payments(r.Context(), username)
	if err != nil {
		statusCode := pkg.ErrorHTTP(err)
		if statusCode == http.StatusInternalServerError {
			h.logger.Err.Printf("could not get data about write-off user %s: %s\n", username, err)
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	if len(info) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	data, errJSON := json.Marshal(&info)
	if errJSON != nil {
		h.logger.Err.Printf("could not marshal json: %s\n", err)
		http.Error(w, errJSON.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(data); err != nil {
		h.logger.Err.Printf("could not write json to response: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
