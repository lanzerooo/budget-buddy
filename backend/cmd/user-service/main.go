package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	// Подключение к PostgreSQL
	db, err := sql.Open("postgres", "user=postgres password=yourpassword dbname=budgetbuddy sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Применение миграций
	if err := goose.Up(db, "internal/user/migrations"); err != nil {
		log.Fatal(err)
	}

	// Инициализация роутера chi
	r := chi.NewRouter()

	// Регистрация маршрута
	r.Post("/api/users/register", registerHandler)

	// Запуск сервера
	log.Println("User Service started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "User registered"}`))
}
