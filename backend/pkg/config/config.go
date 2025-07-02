package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	UserServicePort    string
	FinanceServicePort string
	JWTSecret          string
	DBUrl              string
}

func NewTestConfig() *Config {
    return &Config{
        DBUrl:             "postgres://testuser:testpass@localhost:5433/testdb?sslmode=disable",
        JWTSecret:         "test-secret",
        FinanceServicePort: ":8081",
        UserServicePort:    ":8080",
    }
}


func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &Config{
		UserServicePort:    os.Getenv("USER_SERVICE_PORT"),
		FinanceServicePort: os.Getenv("FINANCE_SERVICE_PORT"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		DBUrl:              os.Getenv("DB_URL"),
	}

	// Проверка обязательных переменных
	if config.UserServicePort == "" {
		return nil, errors.New("USER_SERVICE_PORT environment variable is required")
	}
	if config.FinanceServicePort == "" {
		return nil, errors.New("FINANCE_SERVICE_PORT environment variable is required")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("JWT_SECRET environment variable is required")
	}
	if config.DBUrl == "" {
		return nil, errors.New("DB_URL environment variable is required")
	}

	return config, nil
}
