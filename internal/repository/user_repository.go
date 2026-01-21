package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/davidgrcias/digital-wallet/internal/domain"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByIDForUpdate(ctx context.Context, tx *sql.Tx, id string) (*domain.User, error)
	UpdateBalance(ctx context.Context, tx *sql.Tx, id string, newBalance float64) error
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, name, email, balance, created_at, updated_at 
		FROM users 
		WHERE id = $1
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetByIDForUpdate(ctx context.Context, tx *sql.Tx, id string) (*domain.User, error) {
	query := `
		SELECT id, name, email, balance, created_at, updated_at 
		FROM users 
		WHERE id = $1
		FOR UPDATE
	`

	var user domain.User
	err := tx.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) UpdateBalance(ctx context.Context, tx *sql.Tx, id string, newBalance float64) error {
	query := `
		UPDATE users 
		SET balance = $1, updated_at = NOW() 
		WHERE id = $2
	`

	result, err := tx.ExecContext(ctx, query, newBalance, id)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	n, _ := result.RowsAffected()
	if n == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
