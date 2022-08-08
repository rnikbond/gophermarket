package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gophermarket/internal/repository"
	"gophermarket/internal/repository/postgres"
	"gophermarket/internal/service"
	"gophermarket/internal/tasks"
	"gophermarket/pkg"
	"gophermarket/pkg/logpack"

	market "gophermarket/internal"
	"gophermarket/internal/handlers"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {

	logger := Logger()
	cfg := Config(logger)

	logger.Info.Println(cfg)

	db := PostgresDB(cfg.DatabaseURI, logger)
	pgRepo := PostgresRepository(db, logger)
	loyalty := LoyaltyTask(cfg.AccrualAddress, pgRepo, logger)
	services := service.NewServices(pgRepo, cfg.PasswordSalt, logger)
	handler := handlers.NewHandler(services, cfg.TokenKey, logger)
	server := new(market.Server)

	ctx, cancel := context.WithCancel(context.Background())

	// Запуск сервера
	go func() {
		if err := server.Run(cfg.Address, handler.InitRoutes()); err != nil {
			if err != http.ErrServerClosed {
				logger.Fatal.Fatalf("error occurred while running http server: %s\n", err)
			}
		}
	}()

	// Запуск задачи для обновления статусов и балло из системы лояльности
	if err := loyalty.Scan(ctx); err != nil {
		logger.Fatal.Fatalf("could not run loyalty scanner task: %s\n", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-done

	if err := db.Close(); err != nil {
		logger.Err.Printf("could not close database connection: %s\n", err)
	}

	if err := server.Shutdown(ctx); err != nil {
		if err != http.ErrServerClosed {
			logger.Err.Printf("failed shutdown server: %s\n", err)
		}
	}

	cancel()
}

func Logger() *logpack.LogPack {
	return logpack.NewLogger()
}

func Config(logger *logpack.LogPack) *pkg.Config {
	cfg := pkg.NewConfig()
	if err := cfg.ParseFlags(); err != nil {
		logger.Fatal.Fatalf("could not parse argv: %s\n", err)
	}

	return cfg
}

func PostgresDB(dsn string, logger *logpack.LogPack) *sqlx.DB {

	db, errOpen := sqlx.Open("postgres", dsn)
	if errOpen != nil {
		logger.Fatal.Fatalf("failed to connect to the database: %s\n", errOpen)
	}

	if err := db.Ping(); err != nil {
		logger.Fatal.Fatalf("connection to DB created, but Ping returned error: %s\n", err)
	}

	return db
}

func PostgresRepository(db *sqlx.DB, logger *logpack.LogPack) *repository.Repository {

	repo, errRepo := postgres.NewPostgresRepository(db, logger)
	if errRepo != nil {
		logger.Fatal.Fatalf("failed create repository: %s\n", errRepo)
	}

	return repo
}

func LoyaltyTask(loyaltyAddr string, repo *repository.Repository, logger *logpack.LogPack) tasks.LoyaltyTask {

	return tasks.NewScanner(loyaltyAddr, repo, logger)
}
