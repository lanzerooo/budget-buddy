package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"budgetbuddy/internal/user/models"
	"budgetbuddy/internal/user/repository"
	"budgetbuddy/pkg/auth"
	"budgetbuddy/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

type Handlers struct {
	repo *repository.Repository
}

func NewHandlers(repo *repository.Repository) *Handlers {
	return &Handlers{repo: repo}
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы с фронтенда
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Обрабатываем preflight-запросы
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Передаём управление следующему обработчику
		next(w, r)
	}
}

func SetupRoutes(mux *http.ServeMux, repo *repository.Repository) {
	h := NewHandlers(repo)
	// corsMiddleware к маршрутам
	mux.HandleFunc("/register", corsMiddleware(h.RegisterHandler))
	mux.HandleFunc("/login", corsMiddleware(h.LoginHandler))
}

func (h *Handlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		logger.Error("Failed to decode register request: ", err)
		return
	}

	// Проверка на существование email
	existingUser, err := h.repo.FindUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "Failed to check user existence", http.StatusInternalServerError)
		logger.Error("Failed to check user existence: ", err)
		return
	}
	if existingUser != nil {
		http.Error(w, "Email already registered", http.StatusConflict)
		logger.Error("Email already registered: ", req.Email)
		return
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		logger.Error("Failed to hash password: ", err)
		return
	}

	user := &models.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	// Сохранение пользователя
	_, err = h.repo.SaveUser(user)
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		logger.Error("Failed to save user: ", err)
		return
	}

	// Генерация JWT
	token, err := auth.GenerateJWT(req.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		logger.Error("Failed to generate token: ", err)
		return
	}

	response := models.LoginResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		logger.Error("Failed to decode login request: ", err)
		return
	}

	// Поиск пользователя
	user, err := h.repo.FindUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "Failed to find user", http.StatusInternalServerError)
		logger.Error("Failed to find user: ", err)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		logger.Error("User not found: ", req.Email)
		return
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		logger.Error("Invalid password for user: ", req.Email)
		return
	}

	// Генерация JWT
	token, err := auth.GenerateJWT(req.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		logger.Error("Failed to generate token: ", err)
		return
	}

	response := models.LoginResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
