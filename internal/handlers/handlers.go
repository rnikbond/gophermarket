package handlers

import (
	"net/http"

	"gophermarket/internal/service"
	market "gophermarket/pkg"
	"gophermarket/pkg/logpack"

	"github.com/go-chi/chi"
)

type Handler struct {
	logger   *logpack.LogPack
	services *service.Service
	tokenKey string
}

func NewHandler(services *service.Service, tokenKey string, logger *logpack.LogPack) *Handler {
	return &Handler{
		services: services,
		tokenKey: tokenKey,
		logger:   logger,
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

		// Хендлеры, которые доступны только авторизованным пользователям
		r.Group(func(r chi.Router) {
			r.Use(h.VerifyUser)
			r.Post("/orders", h.CreateOrder)
			r.Get("/orders", h.OrdersList)
			r.Get("/balance", h.Balance)
			r.Post("/balance/withdraw", h.CreateWithPay)
			r.Get("/withdrawals", h.WriteOffInfo)
		})
	})

	return router
}
