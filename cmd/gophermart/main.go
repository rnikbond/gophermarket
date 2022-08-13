package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gophermarket/internal/repository"
	pgPack "gophermarket/internal/repository/postgres"
	"gophermarket/internal/service"
	"gophermarket/internal/tasks"
	"gophermarket/pkg"
	"gophermarket/pkg/logpack"

	market "gophermarket/internal"
	"gophermarket/internal/handlers"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	logger := Logger()
	cfg := Config(logger)

	logger.Info.Println(cfg)

	db := PostgresDB(cfg.DatabaseURI, logger)
	pgRepo := PostgresRepository(db, logger)
	loyalty := LoyaltyTask(cfg.AccrualAddress, pgRepo, cfg.IntervalScan, logger)
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
	loyalty.Scan(ctx)

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

	if err := migrateFor(db.DB, "postgres"); err != nil {
		logger.Fatal.Fatalf("could not apply migrations: %s\n", err)
	}

	return db
}

func LoyaltyTask(loyaltyAddr string, repo *repository.Repository, interval time.Duration, logger *logpack.LogPack) tasks.LoyaltyTask {

	return tasks.NewScanner(loyaltyAddr, repo, interval, logger)
}

func PostgresRepository(db *sqlx.DB, logger *logpack.LogPack) *repository.Repository {

	repo, errRepo := pgPack.NewPostgresRepository(db, logger)
	if errRepo != nil {
		logger.Fatal.Fatalf("failed create repository: %s\n", errRepo)
	}

	return repo
}

func migrateFor(db *sql.DB, driverDB string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		driverDB, driver)

	if err != nil {
		return err
	}

	return m.Up()
}
