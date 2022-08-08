package handlers

import (
	"net/http"

	"gophermarket/pkg"
)

func (h *Handler) VerifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, errCookie := r.Cookie("token")
		if errCookie != nil {
			http.Error(w, pkg.ErrUserUnauthorized.Error(), pkg.ErrorHTTP(pkg.ErrUserUnauthorized))
			return
		}

		bearerToken := cookie.Value
		token, errJWT := h.VerifyJWT(bearerToken)
		if errJWT != nil || !token.Valid {
			http.Error(w, pkg.ErrUserUnauthorized.Error(), pkg.ErrorHTTP(pkg.ErrUserUnauthorized))
			return
		}

		user := token.Claims.(*Token)

		r.SetBasicAuth(user.Username, "")
		next.ServeHTTP(w, r)
	})
}
