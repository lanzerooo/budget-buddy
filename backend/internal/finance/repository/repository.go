package repository

import (
	"budgetbuddy/internal/finance/models"
	"budgetbuddy/pkg/config"
	"budgetbuddy/pkg/logger"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(cfg *config.Config) (*Repository, error) {
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

func (r *Repository) Close() {
	r.db.Close()
}

func (r *Repository) SaveIncome(userID int64, tx *models.Transaction) (int64, error) {
	query := `
        INSERT INTO incomes (user_id, amount, category_id, subcategory_id, description, tags, date, note)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	var id int64
	err := r.db.QueryRow(query, userID, tx.Amount, tx.CategoryID, tx.SubcategoryID, tx.Description, pq.Array(tx.Tags), tx.Date, tx.Note).Scan(&id)
	if err != nil {
		logger.Error("Failed to save income: ", err)
		return 0, err
	}
	return id, nil
}

func (r *Repository) SaveExpense(userID int64, tx *models.Transaction) (int64, error) {
	// Проверка category_id
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS (SELECT 1 FROM categories WHERE id = $1)`, tx.CategoryID).Scan(&exists)
	if err != nil {
		logger.Error("Failed to check category existence: ", err)
		return 0, err
	}
	if !exists {
		logger.Error("Category does not exist: ", tx.CategoryID)
		return 0, fmt.Errorf("category_id %d does not exist", tx.CategoryID)
	}

	// Проверка subcategory_id, если указан
	if tx.SubcategoryID != nil {
		err = r.db.QueryRow(`SELECT EXISTS (SELECT 1 FROM subcategories WHERE id = $1)`, *tx.SubcategoryID).Scan(&exists)
		if err != nil {
			logger.Error("Failed to check subcategory existence: ", err)
			return 0, err
		}
		if !exists {
			logger.Error("Subcategory does not exist: ", *tx.SubcategoryID)
			return 0, fmt.Errorf("subcategory_id %d does not exist", *tx.SubcategoryID)
		}
	}

	query := `
        INSERT INTO expenses (user_id, amount, category_id, subcategory_id, description, tags, date, note)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	var id int64
	err = r.db.QueryRow(query, userID, tx.Amount, tx.CategoryID, tx.SubcategoryID, tx.Description, pq.Array(tx.Tags), tx.Date, tx.Note).Scan(&id)
	if err != nil {
		logger.Error("Failed to save expense: ", err)
		return 0, err
	}
	return id, nil
}

func (r *Repository) GetTransactions(userID int64, txType string) ([]models.Transaction, error) {
	table := "incomes"
	if txType == "expense" {
		table = "expenses"
	}
	query := `
        SELECT id, user_id, amount, category_id, subcategory_id, description, tags, date, note
        FROM ` + table + ` WHERE user_id = $1`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		logger.Error("Failed to get transactions: ", err)
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		var subcategoryID sql.NullInt64
		var tags pq.StringArray
		err := rows.Scan(&tx.ID, &tx.UserID, &tx.Amount, &tx.CategoryID, &subcategoryID, &tx.Description, &tags, &tx.Date, &tx.Note)
		if err != nil {
			logger.Error("Failed to scan transaction: ", err)
			return nil, err
		}
		if subcategoryID.Valid {
			val := subcategoryID.Int64
			tx.SubcategoryID = &val
		}
		tx.Tags = tags
		transactions = append(transactions, tx)
	}
	return transactions, nil
}

func (r *Repository) SaveGoal(userID int64, goal *models.Goal) (int64, error) {
	query := `
        INSERT INTO goals (user_id, name, target_amount, current_amount, deadline, created_at)
        VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var id int64
	err := r.db.QueryRow(query, userID, goal.Name, goal.TargetAmount, goal.CurrentAmount, goal.Deadline, goal.CreatedAt).Scan(&id)
	if err != nil {
		logger.Error("Failed to save goal: ", err)
		return 0, err
	}
	return id, nil
}

func (r *Repository) UpdateGoal(id, userID int64, goal *models.Goal) error {
	query := `
        UPDATE goals SET name=$1, target_amount=$2, current_amount=$3, deadline=$4
        WHERE id=$5 AND user_id=$6`
	_, err := r.db.Exec(query, goal.Name, goal.TargetAmount, goal.CurrentAmount, goal.Deadline, id, userID)
	if err != nil {
		logger.Error("Failed to update goal: ", err)
		return err
	}
	return nil
}

func (r *Repository) GetGoals(userID int64) ([]models.Goal, error) {
	query := `
        SELECT id, user_id, name, target_amount, current_amount, deadline, created_at
        FROM goals WHERE user_id = $1`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		logger.Error("Failed to get goals: ", err)
		return nil, err
	}
	defer rows.Close()

	var goals []models.Goal
	for rows.Next() {
		var g models.Goal
		err := rows.Scan(&g.ID, &g.UserID, &g.Name, &g.TargetAmount, &g.CurrentAmount, &g.Deadline, &g.CreatedAt)
		if err != nil {
			logger.Error("Failed to scan goal: ", err)
			return nil, err
		}
		goals = append(goals, g)
	}
	return goals, nil
}

func (r *Repository) DeleteGoal(id, userID int64) error {
	query := `DELETE FROM goals WHERE id=$1 AND user_id=$2`
	_, err := r.db.Exec(query, id, userID)
	if err != nil {
		logger.Error("Failed to delete goal: ", err)
		return err
	}
	return nil
}

func (r *Repository) SaveCategory(category *models.Category) (int64, error) {
	// Проверка на существование категории
	var exists bool
	err := r.db.QueryRow(`
        SELECT EXISTS (
            SELECT 1 FROM categories WHERE name = $1 AND type = $2
        )`, category.Name, category.Type).Scan(&exists)
	if err != nil {
		logger.Error("Failed to check category existence: ", err)
		return 0, err
	}
	if exists {
		logger.Info("Category already exists: ", category.Name, category.Type)
		var id int64
		err = r.db.QueryRow(`SELECT id FROM categories WHERE name = $1 AND type = $2`, category.Name, category.Type).Scan(&id)
		return id, nil
	}

	query := `INSERT INTO categories (name, type) VALUES ($1, $2) RETURNING id`
	var id int64
	err = r.db.QueryRow(query, category.Name, category.Type).Scan(&id)
	if err != nil {
		logger.Error("Failed to save category: ", err)
		return 0, err
	}
	return id, nil
}

func (r *Repository) SaveSubcategory(subcategory *models.Subcategory) (int64, error) {
	// Проверка на существование category_id
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS (SELECT 1 FROM categories WHERE id = $1)`, subcategory.CategoryID).Scan(&exists)
	if err != nil {
		logger.Error("Failed to check category existence: ", err)
		return 0, err
	}
	if !exists {
		logger.Error("Category does not exist: ", subcategory.CategoryID)
		return 0, fmt.Errorf("category_id %d does not exist", subcategory.CategoryID)
	}

	query := `INSERT INTO subcategories (category_id, name) VALUES ($1, $2) RETURNING id`
	var id int64
	err = r.db.QueryRow(query, subcategory.CategoryID, subcategory.Name).Scan(&id)
	if err != nil {
		logger.Error("Failed to save subcategory: ", err)
		return 0, err
	}
	return id, nil
}

func (r *Repository) GetCategories(userID int64, txType string) ([]models.Category, error) {
	query := `SELECT id, name, type FROM categories WHERE type = $1`
	rows, err := r.db.Query(query, txType)
	if err != nil {
		logger.Error("Failed to get categories: ", err)
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name, &c.Type)
		if err != nil {
			logger.Error("Failed to scan category: ", err)
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *Repository) GetSubcategories(categoryID int64) ([]models.Subcategory, error) {
	query := `SELECT id, category_id, name FROM subcategories WHERE category_id = $1`
	rows, err := r.db.Query(query, categoryID)
	if err != nil {
		logger.Error("Failed to get subcategories: ", err)
		return nil, err
	}
	defer rows.Close()

	var subcategories []models.Subcategory
	for rows.Next() {
		var s models.Subcategory
		err := rows.Scan(&s.ID, &s.CategoryID, &s.Name)
		if err != nil {
			logger.Error("Failed to scan subcategory: ", err)
			return nil, err
		}
		subcategories = append(subcategories, s)
	}
	return subcategories, nil
}

func (r *Repository) SpendingByCategory(userID int64, month string) ([]Spending, error) {
	query := `
        SELECT c.name, SUM(e.amount) as total
        FROM expenses e
        JOIN categories c ON e.category_id = c.id
        WHERE e.user_id = $1 AND TO_CHAR(e.date, 'YYYY-MM') = $2
        GROUP BY c.name`
	rows, err := r.db.Query(query, userID, month)
	if err != nil {
		logger.Error("Failed to get spending data: ", err)
		return nil, err
	}
	defer rows.Close()

	var spending []Spending
	for rows.Next() {
		var s Spending
		err := rows.Scan(&s.Category, &s.Total)
		if err != nil {
			logger.Error("Failed to scan spending data: ", err)
			return nil, err
		}
		spending = append(spending, s)
	}
	return spending, nil
}

type Spending struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

func (r *Repository) IncomeExpenseTrends(userID int64) ([]Trend, error) {
	query := `
        SELECT TO_CHAR(date, 'YYYY-MM') as month, 
               SUM(CASE WHEN t.table_name = 'incomes' THEN amount ELSE 0 END) as income,
               SUM(CASE WHEN t.table_name = 'expenses' THEN amount ELSE 0 END) as expense
        FROM (
            SELECT date, amount, 'incomes' as table_name FROM incomes WHERE user_id = $1
            UNION ALL
            SELECT date, amount, 'expenses' as table_name FROM expenses WHERE user_id = $1
        ) t
        GROUP BY TO_CHAR(date, 'YYYY-MM')
        ORDER BY month`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		logger.Error("Failed to get trends data: ", err)
		return nil, err
	}
	defer rows.Close()

	var trends []Trend
	for rows.Next() {
		var t Trend
		err := rows.Scan(&t.Month, &t.Income, &t.Expense)
		if err != nil {
			logger.Error("Failed to scan trends data: ", err)
			return nil, err
		}
		trends = append(trends, t)
	}
	return trends, nil
}

type Trend struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

func (r *Repository) AverageSpendingByDayOfWeek(userID int64) ([]AverageSpending, error) {
	query := `
        SELECT EXTRACT(DOW FROM date) as day_of_week, AVG(amount) as avg_amount
        FROM expenses
        WHERE user_id = $1
        GROUP BY EXTRACT(DOW FROM date)
        ORDER BY day_of_week`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		logger.Error("Failed to get average spending by day of week: ", err)
		return nil, err
	}
	defer rows.Close()

	var spending []AverageSpending
	for rows.Next() {
		var s AverageSpending
		err := rows.Scan(&s.DayOfWeek, &s.AverageAmount)
		if err != nil {
			logger.Error("Failed to scan average spending data: ", err)
			return nil, err
		}
		spending = append(spending, s)
	}
	return spending, nil
}

type AverageSpending struct {
	DayOfWeek     float64 `json:"day_of_week"` // 0 = Sunday, 1 = Monday, ..., 6 = Saturday
	AverageAmount float64 `json:"average_amount"`
}

func (r *Repository) ForecastSavings(userID, goalID int64) (float64, error) {
	// Проверка на существование цели
	var goal models.Goal
	err := r.db.QueryRow(`
        SELECT target_amount, current_amount
        FROM goals
        WHERE id = $1 AND user_id = $2`, goalID, userID).Scan(&goal.TargetAmount, &goal.CurrentAmount)
	if err == sql.ErrNoRows {
		logger.Error("Goal not found: ", goalID)
		return 0, fmt.Errorf("goal with id %d does not exist for user %d", goalID, userID)
	}
	if err != nil {
		logger.Error("Failed to get goal: ", err)
		return 0, err
	}

	// Рассчитать среднемесячную разницу доходов и расходов
	query := `
        SELECT COALESCE(AVG(income - expense), 0) as avg_savings
        FROM (
            SELECT TO_CHAR(date, 'YYYY-MM') as month,
                   SUM(CASE WHEN t.table_name = 'incomes' THEN amount ELSE 0 END) as income,
                   SUM(CASE WHEN t.table_name = 'expenses' THEN amount ELSE 0 END) as expense
            FROM (
                SELECT date, amount, 'incomes' as table_name FROM incomes WHERE user_id = $1
                UNION ALL
                SELECT date, amount, 'expenses' as table_name FROM expenses WHERE user_id = $1
            ) t
            GROUP BY TO_CHAR(date, 'YYYY-MM')
        ) monthly`
	var avgSavings float64
	err = r.db.QueryRow(query, userID).Scan(&avgSavings)
	if err != nil {
		logger.Error("Failed to calculate average savings: ", err)
		return 0, err
	}

	if avgSavings <= 0 {
		return 0, fmt.Errorf("cannot achieve goal: average savings is zero or negative")
	}

	remainingAmount := goal.TargetAmount - goal.CurrentAmount
	monthsToGoal := remainingAmount / avgSavings
	return monthsToGoal, nil
}
