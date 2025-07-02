package repository

import (
	"budgetbuddy/internal/finance/migrations"
	"budgetbuddy/internal/finance/models"
	"budgetbuddy/pkg/config"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Тесты с sqlmock (юнит-тесты)
func setupTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Failed to create sqlmock")
	return db, mock
}

func TestSaveCategory(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := &Repository{db: db}
	category := &models.Category{Name: "Food", Type: "expense"}

	t.Run("New Category", func(t *testing.T) {
		mock.ExpectQuery(`.*SELECT EXISTS.*`).
			WithArgs("Food", "expense").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		mock.ExpectQuery(`.*INSERT INTO categories.*`).
			WithArgs("Food", "expense").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		id, err := repo.SaveCategory(category)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Existing Category", func(t *testing.T) {
		mock.ExpectQuery(`.*SELECT EXISTS.*`).
			WithArgs("Food", "expense").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		mock.ExpectQuery(`.*SELECT id FROM categories.*`).
			WithArgs("Food", "expense").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		id, err := repo.SaveCategory(category)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DB Error", func(t *testing.T) {
		mock.ExpectQuery(`.*SELECT EXISTS.*`).
			WithArgs("Food", "expense").
			WillReturnError(sql.ErrConnDone)

		id, err := repo.SaveCategory(category)
		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSaveSubcategory(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := &Repository{db: db}
	subcategory := &models.Subcategory{CategoryID: 2, Name: "Groceries"}

	t.Run("Valid Subcategory", func(t *testing.T) {
		mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM categories WHERE id = \$1\)`).
			WithArgs(int64(2)).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		mock.ExpectQuery(`INSERT INTO subcategories \(category_id, name\) VALUES \(\$1, \$2\) RETURNING id`).
			WithArgs(int64(2), "Groceries").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		id, err := repo.SaveSubcategory(subcategory)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Invalid Category", func(t *testing.T) {
		mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM categories WHERE id = \$1\)`).
			WithArgs(int64(2)).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		id, err := repo.SaveSubcategory(subcategory)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "category_id 2 does not exist")
		assert.Equal(t, int64(0), id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSaveExpense(t *testing.T) {
	db, mock := setupTestDB(t)
	defer db.Close()

	repo := &Repository{db: db}
	userID := int64(1)
	tx := &models.Transaction{
		Amount:        200.75,
		CategoryID:    2,
		SubcategoryID: int64Ptr(1),
		Description:   "Grocery shopping",
		Tags:          []string{"food", "expense"},
		Date:          time.Date(2025, 7, 2, 0, 0, 0, 0, time.UTC),
		Note:          "Weekly groceries",
	}

	t.Run("Valid Expense", func(t *testing.T) {
		mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM categories WHERE id = \$1\)`).
			WithArgs(int64(2)).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM subcategories WHERE id = \$1\)`).
			WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		mock.ExpectQuery(`INSERT INTO expenses \(user_id, amount, category_id, subcategory_id, description, tags, date, note\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\) RETURNING id`).
			WithArgs(userID, 200.75, int64(2), int64(1), "Grocery shopping", pq.Array([]string{"food", "expense"}), tx.Date, "Weekly groceries").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		id, err := repo.SaveExpense(userID, tx)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Invalid Category", func(t *testing.T) {
		mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM categories WHERE id = \$1\)`).
			WithArgs(int64(2)).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		id, err := repo.SaveExpense(userID, tx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "category_id 2 does not exist")
		assert.Equal(t, int64(0), id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Invalid Subcategory", func(t *testing.T) {
		mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM categories WHERE id = \$1\)`).
			WithArgs(int64(2)).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		mock.ExpectQuery(`SELECT EXISTS \(SELECT 1 FROM subcategories WHERE id = \$1\)`).
			WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		id, err := repo.SaveExpense(userID, tx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "subcategory_id 1 does not exist")
		assert.Equal(t, int64(0), id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// Вспомогательная функция для указателя на int64
func int64Ptr(i int64) *int64 {
	return &i
}

// Интеграционные тесты с реальной базой
func TestSaveCategoryWithDB(t *testing.T) {
	cfg := config.NewTestConfig()
	db, err := sql.Open("postgres", cfg.DBUrl)
	require.NoError(t, err)
	defer db.Close()

	// Применяем миграции
	err = migrations.RunMigrations(cfg)
	require.NoError(t, err)

	repo := &Repository{db: db}
	category := &models.Category{Name: "Food", Type: "expense"}

	// Очищаем таблицу перед тестом
	_, err = db.Exec("TRUNCATE TABLE categories RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	id, err := repo.SaveCategory(category)
	assert.NoError(t, err)
	assert.NotZero(t, id)

	// Проверяем данные в базе
	var savedName, savedType string
	err = db.QueryRow(`SELECT name, type FROM categories WHERE id = $1`, id).Scan(&savedName, &savedType)
	assert.NoError(t, err)
	assert.Equal(t, "Food", savedName)
	assert.Equal(t, "expense", savedType)
}

func TestSaveSubcategoryWithDB(t *testing.T) {
	cfg := config.NewTestConfig()
	db, err := sql.Open("postgres", cfg.DBUrl)
	require.NoError(t, err)
	defer db.Close()

	// Применяем миграции
	err = migrations.RunMigrations(cfg)
	require.NoError(t, err)

	repo := &Repository{db: db}

	// Создаём категорию
	category := &models.Category{Name: "Food", Type: "expense"}
	catID, err := repo.SaveCategory(category)
	require.NoError(t, err)

	subcategory := &models.Subcategory{CategoryID: catID, Name: "Groceries"}

	// Очищаем таблицу subcategories
	_, err = db.Exec("TRUNCATE TABLE subcategories RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	id, err := repo.SaveSubcategory(subcategory)
	assert.NoError(t, err)
	assert.NotZero(t, id)

	// Проверяем данные в базе
	var savedCatID int64
	var savedName string
	err = db.QueryRow(`SELECT category_id, name FROM subcategories WHERE id = $1`, id).Scan(&savedCatID, &savedName)
	assert.NoError(t, err)
	assert.Equal(t, catID, savedCatID)
	assert.Equal(t, "Groceries", savedName)
}

func TestSaveExpenseWithDB(t *testing.T) {
	cfg := config.NewTestConfig()
	db, err := sql.Open("postgres", cfg.DBUrl)
	require.NoError(t, err)
	defer db.Close()

	// Применяем миграции
	err = migrations.RunMigrations(cfg)
	require.NoError(t, err)

	repo := &Repository{db: db}
	userID := int64(1)

	// Создаём категорию
	category := &models.Category{Name: "Food", Type: "expense"}
	catID, err := repo.SaveCategory(category)
	require.NoError(t, err)

	// Создаём подкатегорию
	subcategory := &models.Subcategory{CategoryID: catID, Name: "Groceries"}
	subcatID, err := repo.SaveSubcategory(subcategory)
	require.NoError(t, err)

	tx := &models.Transaction{
		Amount:        200.75,
		CategoryID:    catID,
		SubcategoryID: int64Ptr(subcatID),
		Description:   "Grocery shopping",
		Tags:          []string{"food", "expense"},
		Date:          time.Date(2025, time.July, 2, 0, 0, 0, 0, time.UTC),
		Note:          "Weekly groceries",
	}

	// Очищаем таблицу expenses
	_, err = db.Exec("TRUNCATE TABLE expenses RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	id, err := repo.SaveExpense(userID, tx)
	assert.NoError(t, err)
	assert.NotZero(t, id)

	// Проверяем данные в базе
	var savedAmount float64
	var savedCatID, savedSubcatID sql.NullInt64
	var savedDesc, savedNote string
	var savedTags []string
	var savedDate time.Time
	err = db.QueryRow(`SELECT amount, category_id, subcategory_id, description, tags, date, note FROM expenses WHERE id = $1`, id).
		Scan(&savedAmount, &savedCatID, &savedSubcatID, &savedDesc, pq.Array(&savedTags), &savedDate, &savedNote)
	assert.NoError(t, err)
	assert.Equal(t, 200.75, savedAmount)
	assert.Equal(t, catID, savedCatID.Int64)
	assert.Equal(t, subcatID, savedSubcatID.Int64)
	assert.Equal(t, "Grocery shopping", savedDesc)
	assert.Equal(t, []string{"food", "expense"}, savedTags)
	assert.Equal(t, tx.Date.UTC(), savedDate.UTC()) // Нормализуем временные зоны
	assert.Equal(t, "Weekly groceries", savedNote)
}
