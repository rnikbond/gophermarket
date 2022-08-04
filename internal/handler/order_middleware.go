package handler

import (
	"fmt"
	"net/http"

	"gophermarket/pkg"
)

func (h *Handler) VerifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
			return
		}

		username, err := h.Username(r)
		if err != nil {
			http.Error(w, err.Error(), pkg.ErrorHTTP(err))
			return
		}

		r.SetBasicAuth(username, "")
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) Username(r *http.Request) (string, error) {

	cookie, err := r.Cookie("token")
	if err != nil {
		fmt.Println("cookie not exists")
		return ``, pkg.ErrUserUnauthorized
	}

	bearerToken := cookie.Value
	token, err := h.VerifyJWT(bearerToken)

	if err != nil || !token.Valid {
		return ``, pkg.ErrUserUnauthorized
	}

	user := token.Claims.(*Token)
	return user.Username, nil
}
