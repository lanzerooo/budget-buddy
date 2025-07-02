package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"budgetbuddy/internal/finance/models"
	finance_repository "budgetbuddy/internal/finance/repository"
	user_repository "budgetbuddy/internal/user/repository"
	"budgetbuddy/pkg/config"
	"budgetbuddy/pkg/logger"

	"github.com/dgrijalva/jwt-go"
)

type Handlers struct {
	repo      *finance_repository.Repository
	userRepo  *user_repository.Repository
	jwtSecret string
}

func NewHandlers(repo *finance_repository.Repository, userRepo *user_repository.Repository, cfg *config.Config) *Handlers {
	return &Handlers{repo: repo, userRepo: userRepo, jwtSecret: cfg.JWTSecret}
}

func SetupRoutes(mux *http.ServeMux, repo *finance_repository.Repository, userRepo *user_repository.Repository, cfg *config.Config) {
	h := NewHandlers(repo, userRepo, cfg)
	mux.HandleFunc("/income", h.authMiddleware(h.AddIncome))
	mux.HandleFunc("/expense", h.authMiddleware(h.AddExpense))
	mux.HandleFunc("/transactions", h.authMiddleware(h.GetTransactions))
	mux.HandleFunc("/categories", h.authMiddleware(h.handleCategories))
	mux.HandleFunc("/subcategories", h.authMiddleware(h.handleSubcategories))
	mux.HandleFunc("/goals", h.authMiddleware(h.handleGoals))
	mux.HandleFunc("/analytics/spending", h.authMiddleware(h.SpendingByCategory))
	mux.HandleFunc("/analytics/trends", h.authMiddleware(h.IncomeExpenseTrends))
	mux.HandleFunc("/analytics/average-spending", h.authMiddleware(h.AverageSpendingByDayOfWeek))
	mux.HandleFunc("/analytics/forecast", h.authMiddleware(h.ForecastSavings))
}

func (h *Handlers) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			logger.Error("Authorization header is empty")
			return
		}
		if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		} else {
			http.Error(w, "Authorization header must start with 'Bearer '", http.StatusUnauthorized)
			logger.Error("Invalid Authorization header format")
			return
		}

		if tokenStr == "" {
			http.Error(w, "JWT token is empty", http.StatusUnauthorized)
			logger.Error("JWT token is empty after removing Bearer prefix")
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.jwtSecret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			logger.Error("Failed to parse or validate token: ", err)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			logger.Error("Failed to parse token claims")
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			http.Error(w, "Invalid email in token", http.StatusUnauthorized)
			logger.Error("Email not found in token claims")
			return
		}

		r.Header.Set("X-User-Email", email)
		next(w, r)
	}
}

func (h *Handlers) AddIncome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		logger.Error("Failed to decode income request: ", err)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Amount must be positive", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
		logger.Error("Invalid date format: ", err)
		return
	}

	userID, err := h.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Failed to get user ID", http.StatusUnauthorized)
		logger.Error("Failed to get user ID: ", err)
		return
	}
	if userID == 0 {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	tx := &models.Transaction{
		UserID:        userID,
		Amount:        req.Amount,
		CategoryID:    req.CategoryID,
		SubcategoryID: req.SubcategoryID,
		Description:   req.Description,
		Tags:          req.Tags,
		Date:          date,
		Note:          req.Note,
	}

	id, err := h.repo.SaveIncome(userID, tx)
	if err != nil {
		http.Error(w, "Failed to save income", http.StatusInternalServerError)
		logger.Error("Failed to save income: ", err)
		return
	}

	response := models.TransactionResponse{
		ID:            id,
		Amount:        tx.Amount,
		CategoryID:    tx.CategoryID,
		SubcategoryID: tx.SubcategoryID,
		Description:   tx.Description,
		Tags:          tx.Tags,
		Date:          tx.Date,
		Note:          tx.Note,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) AddExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		logger.Error("Failed to decode expense request: ", err)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Amount must be positive", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
		logger.Error("Invalid date format: ", err)
		return
	}

	userID, err := h.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Failed to get user ID", http.StatusUnauthorized)
		logger.Error("Failed to get user ID: ", err)
		return
	}
	if userID == 0 {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	tx := &models.Transaction{
		UserID:        userID,
		Amount:        req.Amount,
		CategoryID:    req.CategoryID,
		SubcategoryID: req.SubcategoryID,
		Description:   req.Description,
		Tags:          req.Tags,
		Date:          date,
		Note:          req.Note,
	}

	id, err := h.repo.SaveExpense(userID, tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // Возвращаем текст ошибки
		logger.Error("Failed to save expense: ", err)
		return
	}

	response := models.TransactionResponse{
		ID:            id,
		Amount:        tx.Amount,
		CategoryID:    tx.CategoryID,
		SubcategoryID: tx.SubcategoryID,
		Description:   tx.Description,
		Tags:          tx.Tags,
		Date:          tx.Date,
		Note:          tx.Note,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) GetTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	txType := r.URL.Query().Get("type")
	if txType != "income" && txType != "expense" {
		http.Error(w, "Invalid transaction type, use 'income' or 'expense'", http.StatusBadRequest)
		return
	}

	userID, err := h.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Failed to get user ID", http.StatusUnauthorized)
		logger.Error("Failed to get user ID: ", err)
		return
	}
	if userID == 0 {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	transactions, err := h.repo.GetTransactions(userID, txType)
	if err != nil {
		http.Error(w, "Failed to get transactions", http.StatusInternalServerError)
		logger.Error("Failed to get transactions: ", err)
		return
	}

	response := make([]models.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		response[i] = models.TransactionResponse{
			ID:            tx.ID,
			Amount:        tx.Amount,
			CategoryID:    tx.CategoryID,
			SubcategoryID: tx.SubcategoryID,
			Description:   tx.Description,
			Tags:          tx.Tags,
			Date:          tx.Date,
			Note:          tx.Note,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) handleCategories(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Failed to get user ID", http.StatusUnauthorized)
		logger.Error("Failed to get user ID: ", err)
		return
	}
	if userID == 0 {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodPost {
		var req models.Category
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			logger.Error("Failed to decode category request: ", err)
			return
		}

		if req.Type != "income" && req.Type != "expense" {
			http.Error(w, "Invalid category type, use 'income' or 'expense'", http.StatusBadRequest)
			return
		}

		id, err := h.repo.SaveCategory(&req)
		if err != nil {
			http.Error(w, "Failed to save category", http.StatusInternalServerError)
			logger.Error("Failed to save category: ", err)
			return
		}

		response := models.Category{ID: id, Name: req.Name, Type: req.Type}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
		return
	}

	if r.Method == http.MethodGet {
		txType := r.URL.Query().Get("type")
		if txType != "income" && txType != "expense" {
			http.Error(w, "Invalid transaction type, use 'income' or 'expense'", http.StatusBadRequest)
			return
		}

		categories, err := h.repo.GetCategories(userID, txType)
		if err != nil {
			http.Error(w, "Failed to get categories", http.StatusInternalServerError)
			logger.Error("Failed to get categories: ", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(categories)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h *Handlers) handleSubcategories(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Failed to get user ID", http.StatusUnauthorized)
		logger.Error("Failed to get user ID: ", err)
		return
	}
	if userID == 0 {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodPost {
		var req models.Subcategory
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			logger.Error("Failed to decode subcategory request: ", err)
			return
		}

		id, err := h.repo.SaveSubcategory(&req)
		if err != nil {
			http.Error(w, "Failed to save subcategory", http.StatusInternalServerError)
			logger.Error("Failed to save subcategory: ", err)
			return
		}

		response := models.Subcategory{ID: id, CategoryID: req.CategoryID, Name: req.Name}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
		return
	}

	if r.Method == http.MethodGet {
		categoryIDStr := r.URL.Query().Get("category_id")
		categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid category_id", http.StatusBadRequest)
			return
		}

		subcategories, err := h.repo.GetSubcategories(categoryID)
		if err != nil {
			http.Error(w, "Failed to get subcategories", http.StatusInternalServerError)
			logger.Error("Failed to get subcategories: ", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(subcategories)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h *Handlers) handleGoals(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Failed to get user ID", http.StatusUnauthorized)
		logger.Error("Failed to get user ID: ", err)
		return
	}
	if userID == 0 {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodPost {
		var req models.GoalRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			logger.Error("Failed to decode goal request: ", err)
			return
		}

		deadline, err := time.Parse("2006-01-02", req.Deadline)
		if err != nil {
			http.Error(w, "Invalid deadline format, use YYYY-MM-DD", http.StatusBadRequest)
			logger.Error("Invalid deadline format: ", err)
			return
		}

		goal := &models.Goal{
			UserID:        userID,
			Name:          req.Name,
			TargetAmount:  req.TargetAmount,
			CurrentAmount: 0,
			Deadline:      deadline,
			CreatedAt:     time.Now(),
		}

		id, err := h.repo.SaveGoal(userID, goal)
		if err != nil {
			http.Error(w, "Failed to save goal", http.StatusInternalServerError)
			logger.Error("Failed to save goal: ", err)
			return
		}

		response := models.GoalResponse{
			ID:            id,
			Name:          goal.Name,
			TargetAmount:  goal.TargetAmount,
			CurrentAmount: goal.CurrentAmount,
			Deadline:      goal.Deadline,
			CreatedAt:     goal.CreatedAt,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
		return
	}

	if r.Method == http.MethodPut {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid goal ID", http.StatusBadRequest)
			return
		}

		var req models.GoalRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			logger.Error("Failed to decode goal update request: ", err)
			return
		}

		deadline, err := time.Parse("2006-01-02", req.Deadline)
		if err != nil {
			http.Error(w, "Invalid deadline format, use YYYY-MM-DD", http.StatusBadRequest)
			logger.Error("Invalid deadline format: ", err)
			return
		}

		goal := &models.Goal{
			Name:         req.Name,
			TargetAmount: req.TargetAmount,
			Deadline:     deadline,
		}

		err = h.repo.UpdateGoal(id, userID, goal)
		if err != nil {
			http.Error(w, "Failed to update goal", http.StatusInternalServerError)
			logger.Error("Failed to update goal: ", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodGet {
		goals, err := h.repo.GetGoals(userID)
		if err != nil {
			http.Error(w, "Failed to get goals", http.StatusInternalServerError)
			logger.Error("Failed to get goals: ", err)
			return
		}

		response := make([]models.GoalResponse, len(goals))
		for i, g := range goals {
			response[i] = models.GoalResponse{
				ID:            g.ID,
				Name:          g.Name,
				TargetAmount:  g.TargetAmount,
				CurrentAmount: g.CurrentAmount,
				Deadline:      g.Deadline,
				CreatedAt:     g.CreatedAt,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	if r.Method == http.MethodDelete {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid goal ID", http.StatusBadRequest)
			return
		}

		err = h.repo.DeleteGoal(id, userID)
		if err != nil {
			http.Error(w, "Failed to delete goal", http.StatusInternalServerError)
			logger.Error("Failed to delete goal: ", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h *Handlers) SpendingByCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Failed to get user ID", http.StatusUnauthorized)
		logger.Error("Failed to get user ID: ", err)
		return
	}
	if userID == 0 {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	month := r.URL.Query().Get("month") // Формат: YYYY-MM
	if month == "" {
		http.Error(w, "Month parameter required (YYYY-MM)", http.StatusBadRequest)
		return
	}

	spending, err := h.repo.SpendingByCategory(userID, month)
	if err != nil {
		http.Error(w, "Failed to get spending data", http.StatusInternalServerError)
		logger.Error("Failed to get spending data: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(spending)
}

func (h *Handlers) IncomeExpenseTrends(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Failed to get user ID", http.StatusUnauthorized)
		logger.Error("Failed to get user ID: ", err)
		return
	}
	if userID == 0 {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	trends, err := h.repo.IncomeExpenseTrends(userID)
	if err != nil {
		http.Error(w, "Failed to get trends data", http.StatusInternalServerError)
		logger.Error("Failed to get trends data: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trends)
}

func (h *Handlers) getUserIDFromToken(r *http.Request) (int64, error) {
	email := r.Header.Get("X-User-Email")
	return h.userRepo.GetUserIDByEmail(email)
}

func (h *Handlers) AverageSpendingByDayOfWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Failed to get user ID", http.StatusUnauthorized)
		logger.Error("Failed to get user ID: ", err)
		return
	}
	if userID == 0 {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	spending, err := h.repo.AverageSpendingByDayOfWeek(userID)
	if err != nil {
		http.Error(w, "Failed to get average spending data", http.StatusInternalServerError)
		logger.Error("Failed to get average spending data: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(spending)
}

func (h *Handlers) ForecastSavings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Failed to get user ID", http.StatusUnauthorized)
		logger.Error("Failed to get user ID: ", err)
		return
	}
	if userID == 0 {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	goalIDStr := r.URL.Query().Get("goal_id")
	goalID, err := strconv.ParseInt(goalIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid goal ID", http.StatusBadRequest)
		return
	}

	monthsToGoal, err := h.repo.ForecastSavings(userID, goalID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // Возвращаем текст ошибки
		logger.Error("Failed to forecast savings: ", err)
		return
	}

	response := map[string]float64{"months_to_goal": monthsToGoal}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
