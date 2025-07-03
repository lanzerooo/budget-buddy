package models

import "time"

type TransactionRequest struct {
	Amount        float64  `json:"amount"`
	CategoryID    int64    `json:"category_id"`
	SubcategoryID *int64   `json:"subcategory_id,omitempty"`
	Description   string   `json:"description"`
	Tags          []string `json:"tags,omitempty"`
	Date          string   `json:"date"`
	Note          string   `json:"note"`
}

type Transaction struct {
	ID            int64
	UserID        int64
	Amount        float64
	CategoryID    int64
	SubcategoryID *int64
	Description   string
	Tags          []string
	Date          time.Time
	Note          string
}

type TransactionResponse struct {
	ID            int64     `json:"id"`
	Amount        float64   `json:"amount"`
	CategoryID    int64     `json:"category_id"`
	SubcategoryID *int64    `json:"subcategory_id,omitempty"`
	Description   string    `json:"description"`
	Tags          []string  `json:"tags,omitempty"`
	Date          time.Time `json:"date"`
	Note          string    `json:"note"`
}

type GoalRequest struct {
	Name         string  `json:"name"`
	TargetAmount float64 `json:"target_amount"`
	Deadline     string  `json:"deadline"`
}

type Goal struct {
	ID            int64
	UserID        int64
	Name          string
	TargetAmount  float64
	CurrentAmount float64
	Deadline      time.Time
	CreatedAt     time.Time
}

type GoalResponse struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	TargetAmount  float64   `json:"target_amount"`
	CurrentAmount float64   `json:"current_amount"`
	Deadline      time.Time `json:"deadline"`
	CreatedAt     time.Time `json:"created_at"`
}

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Subcategory struct {
	ID         int64  `json:"id"`
	CategoryID int64  `json:"category_id"`
	Name       string `json:"name"`
}

type Budget struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	CategoryID int64     `json:"category_id"`
	Amount     float64   `json:"amount"`
	Month      string    `json:"month"`
	CreatedAt  time.Time `json:"created_at"`
}
