package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	market "gophermarket/internal"
	"gophermarket/pkg"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			h.logger.Err.Printf("could not close request body: %s\n", err)
		}
	}()

	var user market.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("SignUp :: Error unmarshal %s to User struct\n", r.Body)
		return
	}

	errSignUp := h.services.Auth.SignUp(user)
	if errSignUp != nil {
		log.Printf("SignUp :: service return error: %v\n", errSignUp)
		http.Error(w, errSignUp.Error(), pkg.ErrorHTTP(errSignUp))
		return
	}

	if err := saveAuth(&w, h.GenerateJWT(user)); err != nil {
		log.Printf("SignUp :: saveAuth return error: %v\n", errSignUp)
		http.Error(w, err.Error(), pkg.ErrorHTTP(err))
		return
	}
}
