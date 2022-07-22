package server

import (
	"fmt"
	"net/http"

	"gophermarket/internal/server/handlers"

	"github.com/go-chi/chi"
)

func Run() *http.Server {

	router := chi.NewRouter()
	router.With()
	router.Post("/api/user/login", handlers.Login())
	router.Post("/api/user/register", handlers.Register())

	serverHTTP := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: router,
	}

	go func() {
		if err := serverHTTP.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("HTTP server ListenAndServe: %v\n", err)
		}
	}()

	return serverHTTP
}
