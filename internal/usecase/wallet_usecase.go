package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/davidgrcias/digital-wallet/internal/domain"
	"github.com/davidgrcias/digital-wallet/internal/repository"
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

	// Log attempt
	log.Printf("[Withdraw] Processing request for user: %s, amount: %.2f", userID, req.Amount)

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[Withdraw] Failed to begin tx: %v", err)
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	user, err := u.userRepo.GetByIDForUpdate(ctx, tx, userID)
	if err != nil {
		tx.Rollback()
		log.Printf("[Withdraw] Failed to lock user %s: %v", userID, err)
		return nil, err
	}

	if user.Balance < req.Amount {
		tx.Rollback()
		log.Printf("[Withdraw] Insufficient balance for user %s (Balance: %.2f, Req: %.2f)", userID, user.Balance, req.Amount)
		return nil, domain.ErrInsufficientBalance
	}

	balanceBefore := user.Balance
	balanceAfter := user.Balance - req.Amount

	if err = u.userRepo.UpdateBalance(ctx, tx, userID, balanceAfter); err != nil {
		tx.Rollback()
		log.Printf("[Withdraw] Failed to update balance: %v", err)
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
		log.Printf("[Withdraw] Failed to create transaction record: %v", err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[Withdraw] Failed to commit: %v", err)
		return nil, fmt.Errorf("commit: %w", err)
	}

	log.Printf("[Withdraw] Success! User: %s, TxID: %s, NewBalance: %.2f", userID, saved.ID, balanceAfter)

	return &domain.WithdrawResponse{
		TransactionID: saved.ID,
		UserID:        userID,
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		CreatedAt:     saved.CreatedAt,
	}, nil
}
