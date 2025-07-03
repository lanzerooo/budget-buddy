package repository

import (
	"database/sql"

	"budgetbuddy/internal/user/models"
	"budgetbuddy/pkg/config"
	"budgetbuddy/pkg/logger"

	_ "github.com/lib/pq"
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

func (r *Repository) SaveUser(user *models.User) (int64, error) {
	query := `INSERT INTO users (email, password, name, created_at) 
              VALUES ($1, $2, $3, $4) RETURNING id`
	var id int64
	err := r.db.QueryRow(query, user.Email, user.Password, user.Name, user.CreatedAt).Scan(&id)
	if err != nil {
		logger.Error("Failed to save user: ", err)
		return 0, err
	}
	return id, nil
}

func (r *Repository) FindUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password, name, created_at FROM users WHERE email = $1`
	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logger.Error("Failed to find user: ", err)
		return nil, err
	}
	return user, nil
}

func (r *Repository) GetUserIDByEmail(email string) (int64, error) {
	query := `SELECT id FROM users WHERE email = $1`
	var id int64
	err := r.db.QueryRow(query, email).Scan(&id)
	if err == sql.ErrNoRows {
		logger.Error("User not found for email: ", email)
		return 0, nil
	}
	if err != nil {
		logger.Error("Failed to get user ID: ", err)
		return 0, err
	}
	return id, nil
}

//добавлен метод для профиля и пароля
func (r *Repository) GetUserProfile(userID int64) (*models.User, error) {
	query := `SELECT id, email, name, created_at FROM users WHERE id = $1`
	user := &models.User{}
	err := r.db.QueryRow(query, userID).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt)
	if err == sql.ErrNoRows {
		logger.Error("User not found: ", userID)
		return nil, nil
	}
	if err != nil {
		logger.Error("Failed to get user profile: ", err)
		return nil, err
	}
	return user, nil
}

func (r *Repository) UpdateUserName(userID int64, name string) error {
	query := `UPDATE users SET name = $1 WHERE id = $2`
	_, err := r.db.Exec(query, name, userID)
	if err != nil {
		logger.Error("Failed to update user name: ", err)
		return err
	}
	return nil
}

func (r *Repository) UpdateUserPassword(userID int64, hashedPassword string) error {
	query := `UPDATE users SET password = $1 WHERE id = $2`
	_, err := r.db.Exec(query, hashedPassword, userID)
	if err != nil {
		logger.Error("Failed to update user password: ", err)
		return err
	}
	return nil
}
