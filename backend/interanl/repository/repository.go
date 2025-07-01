package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lanzerooo/budget-buddy.git/budjet-buddy/interanl/model"
)

type Repository interface {
	CreateTransaction(tx model.Transaction) error
	GetAllTransactions() ([]model.Transaction, error)
	GetBalance() (float64, error)
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepository (db *pgxpool.Pool) Repository {
	return &repo{db: db}
}

func (r *repo) CreateTransaction(tx model.Transaction) error {
	_, err := r.db.Exec(context.Background(),
	    "INSERT INTO transactions (amount, category, created_at) VALUES ($1, $2, $3)",
		tx.Amount, tx.Category, tx.CreatedAt)
	return err
}

func (r *repo) GetAllTransactions() ([]model.Transaction, error) {
	rows, err := r.db.Query(context.Background(), "SELECT id, amount, category, created_at FROM transactions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Transaction
	for rows.Next() {
		tx := model.Transaction{}
		err := rows.Scan(&tx.ID, &tx.Amount, &tx.Category, &tx.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, tx)
	}
	return result, nil
}

func (r *repo) GetBalance() (float64, error) {
	var total float64
	err := r.db.QueryRow(context.Background(), "SELECT COALESCE(SUM(amount),0) FROM transactions").Scan(&total)
	return total, err
}