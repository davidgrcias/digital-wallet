package domain

import (
	"errors"
	"time"
)

const (
	TransactionTypeWithdraw = "WITHDRAW"
	TransactionTypeDeposit  = "DEPOSIT"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidAmount       = errors.New("amount must be greater than 0")
	ErrUserNotFound        = errors.New("user not found")
)

type Transaction struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
}

type WithdrawRequest struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

func (w *WithdrawRequest) Validate() error {
	if w.Amount <= 0 {
		return ErrInvalidAmount
	}
	return nil
}

type WithdrawResponse struct {
	TransactionID string    `json:"transaction_id"`
	UserID        string    `json:"user_id"`
	Amount        float64   `json:"amount"`
	BalanceBefore float64   `json:"balance_before"`
	BalanceAfter  float64   `json:"balance_after"`
	CreatedAt     time.Time `json:"created_at"`
}
