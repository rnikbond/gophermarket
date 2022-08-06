package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gophermarket/internal/repository/postgres"
	"gophermarket/internal/service"
	"gophermarket/internal/tasks"
	"gophermarket/pkg"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	market "gophermarket/internal"
	"gophermarket/internal/handler"
)

func configureLog() {

	logrus.SetFormatter(new(logrus.JSONFormatter))
}

func pgDriver(dsn string) (*sqlx.DB, error) {

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {

	configureLog()

	cfg := pkg.NewConfig()
	if err := cfg.ParseFlags(); err != nil {
		logrus.Fatalf("error read argv: %v\n", err)
	}

	fmt.Println(cfg)

	db, errDB := pgDriver(cfg.DatabaseURI)
	if errDB != nil {
		logrus.Fatalf("error create connection to database: %v\n", errDB)
	}

	pgRepo, errRepo := postgres.NewPostgresRepository(db)
	if errRepo != nil {
		logrus.Fatalf("error create postgres repository: %v\n", errRepo)
	}

	accrualScanner := tasks.NewScanner(cfg.AccrualAddress, pgRepo)
	services := service.NewServices(pgRepo, cfg.PasswordSalt)
	handlers := handler.NewHandler(services, cfg.TokenKey)
	srv := new(market.Server)

	go func() {
		if err := srv.Run(cfg.Address, handlers.InitRoutes()); err != nil {
			if err != http.ErrServerClosed {
				logrus.Fatalf("error occured while running http server: %v\n", err)
			}
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	if err := accrualScanner.Scan(ctx); err != nil {
		logrus.Fatalf("error run accrual scanner task: %v\n", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-done

	if err := db.Close(); err != nil {
		logrus.Fatalf("error close database connection:%+v", err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		if err != http.ErrServerClosed {
			logrus.Fatalf("Server Shutdown Failed:%+v", err)
		}
	}

	cancel()
}
