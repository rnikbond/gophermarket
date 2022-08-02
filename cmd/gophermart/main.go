package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gophermarket/pkg/repository/postgres"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"gophermarket/config"
	market "gophermarket/pkg"
	"gophermarket/pkg/handler"
	"gophermarket/pkg/service"
)

func init() {

	logrus.SetFormatter(new(logrus.JSONFormatter))
}

func main() {

	cfg := config.NewConfig()
	if err := cfg.ParseFlags(); err != nil {
		logrus.Fatalf("error read argv: %v\n", err)
	}

	log.Println(cfg)

	pgStorage, errDB := postgres.NewPostgresRepository(cfg.DatabaseURI)
	if errDB != nil {
		logrus.Fatalf("error create database storage: %v\n", errDB)
	}

	services := service.NewService(pgStorage)
	handlers := handler.NewHandler(services)
	srv := new(market.Server)

	go func() {
		if err := srv.Run(":8080", handlers.InitRoutes()); err != nil {
			if err != http.ErrServerClosed {
				logrus.Fatalf("error occured while running http server: %v\n", err)
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

	if err := pgStorage.Finish(); err != nil {
		logrus.Fatalf("error close database connection:%+v", err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		if err != http.ErrServerClosed {
			logrus.Fatalf("Server Shutdown Failed:%+v", err)
		}
	}
}
