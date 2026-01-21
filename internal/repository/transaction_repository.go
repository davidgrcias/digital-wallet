package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/davidgrcias/digital-wallet/internal/domain"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *sql.Tx, transaction *domain.Transaction) (*domain.Transaction, error)
}

type PostgresTransactionRepository struct {
	db *sql.DB
}

func NewPostgresTransactionRepository(db *sql.DB) TransactionRepository {
	return &PostgresTransactionRepository{db: db}
}

func (r *PostgresTransactionRepository) Create(ctx context.Context, tx *sql.Tx, transaction *domain.Transaction) (*domain.Transaction, error) {
	query := `
		INSERT INTO transactions (user_id, type, amount, balance_before, balance_after, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, type, amount, balance_before, balance_after, description, created_at
	`

	var result domain.Transaction
	var description sql.NullString

	err := tx.QueryRowContext(
		ctx,
		query,
		transaction.UserID,
		transaction.Type,
		transaction.Amount,
		transaction.BalanceBefore,
		transaction.BalanceAfter,
		nullString(transaction.Description),
	).Scan(
		&result.ID,
		&result.UserID,
		&result.Type,
		&result.Amount,
		&result.BalanceBefore,
		&result.BalanceAfter,
		&description,
		&result.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("insert failed: %w", err)
	}

	if description.Valid {
		result.Description = description.String
	}

	return &result, nil
}

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
