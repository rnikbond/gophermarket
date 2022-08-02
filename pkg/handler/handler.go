package handler

import (
	"net/http"

	market "gophermarket/pkg"
	"gophermarket/pkg/service"

	"github.com/go-chi/chi"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func SetCookie(w *http.ResponseWriter, token string) error {

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

	router.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/login", h.SignIn)
			r.Post("/register", h.SignUp)
			r.Post("/orders", h.CreateOrder)
		})
	})

	return router
}
