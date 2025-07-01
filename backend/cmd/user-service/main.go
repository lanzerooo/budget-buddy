package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"budgetbuddy/internal/user/handlers"
	"budgetbuddy/internal/user/migrations"
	"budgetbuddy/internal/user/repository"
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

	// Инициализация репозитория
	repo, err := repository.NewRepository(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize repository: ", err)
	}
	defer repo.Close()

	// Инициализация роутера
	mux := http.NewServeMux()

	// Инициализация обработчиков
	handlers.SetupRoutes(mux, repo)

	// Настройка сервера
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Запуск сервера в горутине
	go func() {
		logger.Info("Starting server on port ", cfg.ServerPort)
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
	logger.Info("Server gracefully stopped")
}
