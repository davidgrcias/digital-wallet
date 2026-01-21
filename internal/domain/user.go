package domain

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BalanceResponse struct {
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Balance     float64   `json:"balance"`
	LastUpdated time.Time `json:"last_updated"`
}
