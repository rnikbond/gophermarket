package handler

import (
	"net/http"

	"gophermarket/internal/service"
	market "gophermarket/pkg"

	"github.com/go-chi/chi"
)

type Handler struct {
	services *service.Service
	tokenKey string
}

func NewHandler(services *service.Service, tokenKey string) *Handler {
	return &Handler{
		services: services,
		tokenKey: tokenKey,
	}
}

func saveAuth(w *http.ResponseWriter, token string) error {

	if len(token) < 1 {
		return market.ErrGenerateToken
	}

	http.SetCookie(*w, &http.Cookie{
		Name:  "token",
		Value: token,
	})

	return nil
}

func (h *Handler) InitRoutes() *chi.Mux {

	router := chi.NewRouter()

	router.Route("/api/user", func(r chi.Router) {
		r.Post("/login", h.SignIn)
		r.Post("/register", h.SignUp)
		r.Group(func(r chi.Router) {
			r.Use(h.VerifyUser)
			r.Post("/orders", h.CreateOrder)
			r.Get("/balance", h.Balance)
		})
	})

	return router
}
