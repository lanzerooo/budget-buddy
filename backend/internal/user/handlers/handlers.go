package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"budgetbuddy/internal/user/models"
	"budgetbuddy/internal/user/repository"
	"budgetbuddy/pkg/auth"
	"budgetbuddy/pkg/config"
	"budgetbuddy/pkg/logger"
	"budgetbuddy/pkg/middleware"

	"golang.org/x/crypto/bcrypt"
)

type Handlers struct {
	repo      *repository.Repository
	jwtSecret string
}

func NewHandlers(repo *repository.Repository, cfg *config.Config) *Handlers {
	return &Handlers{
		repo:      repo,
		jwtSecret: cfg.JWTSecret,
	}
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

func SetupRoutes(mux *http.ServeMux, repo *repository.Repository, cfg *config.Config) {
	h := NewHandlers(repo, cfg)
	// corsMiddleware к маршрутам
	mux.HandleFunc("/register", corsMiddleware(h.RegisterHandler))
	mux.HandleFunc("/login", corsMiddleware(h.LoginHandler))
	mux.HandleFunc("/profile", corsMiddleware(middleware.AuthMiddleware(h.jwtSecret, h.GetProfile)))
	mux.HandleFunc("/profile/update", corsMiddleware(middleware.AuthMiddleware(h.jwtSecret, h.UpdateProfile)))
	mux.HandleFunc("/password", corsMiddleware(middleware.AuthMiddleware(h.jwtSecret, h.UpdatePassword)))
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

	_, err = h.repo.SaveUser(user)
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		logger.Error("Failed to save user: ", err)
		return
	}

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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		logger.Error("Invalid password for user: ", req.Email)
		return
	}

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

//добавлены новые хэндлеры

func (h *Handlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, err := h.getUserIDFromToken(r)
	if err != nil || userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		logger.Error("Failed to get user ID: ", err)
		return
	}
	user, err := h.repo.GetUserProfile(userID)
	if err != nil {
		http.Error(w, "Failed to get profile", http.StatusInternalServerError)
		logger.Error("Failed to get profile: ", err)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	response := models.UserProfileResponse{
		Email: user.Email,
		Name:  user.Name,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, err := h.getUserIDFromToken(r)
	if err != nil || userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		logger.Error("Failed to get user ID: ", err)
		return
	}
	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		logger.Error("Failed to decode update profile request: ", err)
		return
	}
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	if err := h.repo.UpdateUserName(userID, req.Name); err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		logger.Error("Failed to update profile: ", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, err := h.getUserIDFromToken(r)
	if err != nil || userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		logger.Error("Failed to get user ID: ", err)
		return
	}
	var req models.UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		logger.Error("Failed to decode update password request: ", err)
		return
	}
	if req.NewPassword == "" || req.OldPassword == "" {
		http.Error(w, "Old and new passwords are required", http.StatusBadRequest)
		return
	}
	user, err := h.repo.FindUserByEmail(r.Header.Get("X-User-Email"))
	if err != nil || user == nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		logger.Error("User not found: ", r.Header.Get("X-User-Email"))
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		http.Error(w, "Invalid old password", http.StatusUnauthorized)
		logger.Error("Invalid old password for user: ", r.Header.Get("X-User-Email"))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		logger.Error("Failed to hash password: ", err)
		return
	}
	if err := h.repo.UpdateUserPassword(userID, string(hashedPassword)); err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		logger.Error("Failed to update password: ", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) getUserIDFromToken(r *http.Request) (int64, error) {
	return h.repo.GetUserIDByEmail(r.Header.Get("X-User-Email"))
}
