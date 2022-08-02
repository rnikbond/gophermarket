package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/caarlos0/env"
)

type Config struct {
	Address        string `env:"RUN_ADDRESS "`
	DatabaseURI    string `env:"DATABASE_URI"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func NewConfig() Config {
	cfg := Config{
		Address: ":8080",
	}

	cfg.ReadEnvVars()
	return cfg
}

func (cfg Config) String() string {

	s := "\n"
	s += fmt.Sprintf("RUN_ADDRESS: %s\n", cfg.Address)
	s += fmt.Sprintf("DATABASE_URI: %s\n", cfg.DatabaseURI)
	s += fmt.Sprintf("ACCRUAL_SYSTEM_ADDRESS: %s\n", cfg.AccrualAddress)

	return s
}

// ReadEnvVars - Чтение переменных среды
func (cfg *Config) ReadEnvVars() {

	// Чтение переменных среды
	if err := env.Parse(cfg); err != nil {
		log.Println(err)
	}

	// Убираем пробелы из адреса
	cfg.Address = strings.TrimSpace(cfg.Address)
}

// ParseFlags - Разбор аргументов командной строки
func (cfg *Config) ParseFlags() error {

	flag.StringVar(&cfg.AccrualAddress, "r", cfg.AccrualAddress, "string - accrual address")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "string - database DSN")
	addr := flag.String("a", cfg.Address, "string - host:port")
	flag.Parse()

	if addr == nil || *addr == "" {
		return nil
	}

	parsedAddr := strings.Split(*addr, ":")
	if len(parsedAddr) != 2 {
		return errors.New("invalid address format")
	}

	if len(parsedAddr[0]) > 0 && parsedAddr[0] != "localhost" {
		if ip := net.ParseIP(parsedAddr[0]); ip == nil {
			return errors.New("incorrect ip: " + parsedAddr[0])
		}
	}

	if _, err := strconv.Atoi(parsedAddr[1]); err != nil {
		return errors.New("incorrect port: " + parsedAddr[1])
	}

	cfg.Address = *addr
	return nil
}
