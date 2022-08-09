package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	market "gophermarket/internal"
	"gophermarket/pkg"
)

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("Error :: Handler 'SignIn' :: error close body: %v\n", err)
		}
	}()

	var user market.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	errSignIn := h.services.Auth.SignIn(r.Context(), user)
	if errSignIn != nil {
		http.Error(w, errSignIn.Error(), pkg.ErrorHTTP(errSignIn))
		return
	}

	if err := saveAuth(&w, h.GenerateJWT(user)); err != nil {
		http.Error(w, err.Error(), pkg.ErrorHTTP(err))
		return
	}
}
