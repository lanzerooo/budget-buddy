package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"budgetbuddy/internal/finance/handlers"
	"budgetbuddy/internal/finance/migrations"
	finance_repository "budgetbuddy/internal/finance/repository"
	user_repository "budgetbuddy/internal/user/repository"
	"budgetbuddy/pkg/config"
	"budgetbuddy/pkg/logger"
)

func main() {
	// Инициализация логгера
	logger.Init()

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config: ", err)
	}

	// Выполнение миграций
	if err := migrations.RunMigrations(cfg); err != nil {
		logger.Fatal("Failed to run migrations: ", err)
	}

	// Инициализация репозитория Finance Service
	repo, err := finance_repository.NewRepository(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize finance repository: ", err)
	}
	defer repo.Close()

	// Инициализация репозитория User Service
	userRepo, err := user_repository.NewRepository(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize user repository: ", err)
	}
	defer userRepo.Close()

	// Инициализация роутера
	mux := http.NewServeMux()

	// Инициализация обработчиков
	handlers.SetupRoutes(mux, repo, userRepo, cfg)

	// Настройка сервера
	server := &http.Server{
		Addr:         ":" + cfg.FinanceServicePort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Запуск сервера в горутине
	go func() {
		logger.Info("Starting finance service on port ", cfg.FinanceServicePort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed: ", err)
		}
	}()

	// Ожидание сигнала завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	// Корректное завершение работы
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server shutdown failed: ", err)
	}
	logger.Info("Finance service gracefully stopped")
}
