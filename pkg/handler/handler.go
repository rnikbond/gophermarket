package handler

import (
	"gophermarket/pkg/service"

	"github.com/go-chi/chi"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *chi.Mux {

	router := chi.NewRouter()

	router.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/login", h.login)
			r.Post("/register", h.register)
		})
	})

	return router
}
