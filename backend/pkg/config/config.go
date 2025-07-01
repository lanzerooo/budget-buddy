package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	JWTSecret  string
	DBUrl      string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &Config{
		ServerPort: os.Getenv("PORT"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		DBUrl:      os.Getenv("DB_URL"),
	}

	// Проверка обязательных переменных
	if config.ServerPort == "" {
		return nil, errors.New("PORT environment variable is required")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("JWT_SECRET environment variable is required")
	}
	if config.DBUrl == "" {
		return nil, errors.New("DB_URL environment variable is required")
	}

	return config, nil
}
