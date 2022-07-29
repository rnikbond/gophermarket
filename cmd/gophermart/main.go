package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gophermarket/config"
	gophermarket "gophermarket/pkg"
	"gophermarket/pkg/handler"
	"gophermarket/pkg/repository"
	"gophermarket/pkg/service"
)

func main() {

	cfg := config.NewConfig()
	if err := cfg.ParseFlags(); err != nil {
		log.Fatalf("error read argv: %v\n", err)
	}

	log.Println(cfg)

	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(gophermarket.Server)

	go func() {
		if err := srv.Run(":8080", handlers.InitRoutes()); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalf("error occured while running http server: %v\n", err)
			}
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		if err != http.ErrServerClosed {
			log.Fatalf("Server Shutdown Failed:%+v", err)
		}
	}
}
