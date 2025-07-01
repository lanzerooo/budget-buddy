package model

import "time"

type Transaction struct {
	ID        int     `json:"id"`
	Amount    float64 `json:"amount"`
	Category  string  `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}