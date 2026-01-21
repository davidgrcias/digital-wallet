package usecase

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/davidgarcia/digital-wallet/internal/domain"
	"github.com/davidgarcia/digital-wallet/internal/repository"
)

type WalletUsecase interface {
	GetBalance(ctx context.Context, userID string) (*domain.BalanceResponse, error)
	Withdraw(ctx context.Context, userID string, req *domain.WithdrawRequest) (*domain.WithdrawResponse, error)
}

type WalletUsecaseImpl struct {
	db              *sql.DB
	userRepo        repository.UserRepository
	transactionRepo repository.TransactionRepository
}

func NewWalletUsecase(
	db *sql.DB,
	userRepo repository.UserRepository,
	transactionRepo repository.TransactionRepository,
) WalletUsecase {
	return &WalletUsecaseImpl{
		db:              db,
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
	}
}

func (u *WalletUsecaseImpl) GetBalance(ctx context.Context, userID string) (*domain.BalanceResponse, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	return &domain.BalanceResponse{
		UserID:      user.ID,
		Name:        user.Name,
		Balance:     user.Balance,
		LastUpdated: user.UpdatedAt,
	}, nil
}

func (u *WalletUsecaseImpl) Withdraw(ctx context.Context, userID string, req *domain.WithdrawRequest) (*domain.WithdrawResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	user, err := u.userRepo.GetByIDForUpdate(ctx, tx, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if user.Balance < req.Amount {
		tx.Rollback()
		return nil, domain.ErrInsufficientBalance
	}

	balanceBefore := user.Balance
	balanceAfter := user.Balance - req.Amount

	if err = u.userRepo.UpdateBalance(ctx, tx, userID, balanceAfter); err != nil {
		tx.Rollback()
		return nil, err
	}

	txn := &domain.Transaction{
		UserID:        userID,
		Type:          domain.TransactionTypeWithdraw,
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Description:   req.Description,
	}

	saved, err := u.transactionRepo.Create(ctx, tx, txn)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &domain.WithdrawResponse{
		TransactionID: saved.ID,
		UserID:        userID,
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		CreatedAt:     saved.CreatedAt,
	}, nil
}
