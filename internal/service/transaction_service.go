package service

import (
	"context"
	"errors"
	"fmt"
	"internal-transfer-system/internal/database"
	"internal-transfer-system/internal/models"
	"internal-transfer-system/internal/repository"
	"math/big"
	"time"

	"gorm.io/gorm"
)

type TransactionService struct {
	accountRepo     *repository.AccountRepository
	transactionRepo *repository.TransactionRepository
	db              *gorm.DB
}

func NewTransactionService() *TransactionService {
	return &TransactionService{
		accountRepo:     repository.NewAccountRepository(database.GetDB()),
		transactionRepo: repository.NewTransactionRepository(database.GetDB()),
		db:              database.GetDB(),
	}
}

// CreateTransfer processes a transfer transaction between two accounts
// Uses database transaction to ensure atomicity - if any step fails, all changes are rolled back automatically
func (s *TransactionService) CreateTransfer(ctx context.Context, req *models.CreateTransactionRequest) error {
	// Validate amount
	amount, ok := new(big.Float).SetString(req.Amount)
	if !ok {
		return fmt.Errorf("invalid amount format")
	}

	// Use database transaction to ensure atomicity
	// GORM's Transaction method automatically:
	// - Commits if the function returns nil
	// - Rolls back if the function returns an error
	// - Rolls back on panic
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create temporary repositories with transaction context
		accountRepo := repository.NewAccountRepository(tx)
		transactionRepo := repository.NewTransactionRepository(tx)

		// Step 1: Get source account balance with row-level lock (SELECT FOR UPDATE)
		// If this fails, transaction will rollback
		sourceBalanceStr, err := accountRepo.GetBalanceForUpdate(ctx, tx, req.SourceAccountID)
		if err != nil {
			if errors.Is(err, repository.ErrAccountNotFound) {
				return fmt.Errorf("source account not found")
			}
			return fmt.Errorf("failed to get source account balance: %w", err)
		}

		// Step 2: Validate source account balance format
		// If this fails, transaction will rollback
		sourceBalance, ok := new(big.Float).SetString(sourceBalanceStr)
		if !ok {
			return fmt.Errorf("invalid source account balance format")
		}

		// Step 3: Check if source account has sufficient balance
		// If insufficient, transaction will rollback
		if sourceBalance.Cmp(amount) < 0 {
			return repository.ErrInsufficientBalance
		}

		// Step 4: Get destination account with row-level lock
		// If this fails, transaction will rollback
		destBalanceStr, err := accountRepo.GetBalanceForUpdate(ctx, tx, req.DestinationAccountID)
		if err != nil {
			if errors.Is(err, repository.ErrAccountNotFound) {
				return fmt.Errorf("destination account not found")
			}
			return fmt.Errorf("failed to get destination account balance: %w", err)
		}

		// Step 5: Validate destination account balance format
		// If this fails, transaction will rollback
		destBalance, ok := new(big.Float).SetString(destBalanceStr)
		if !ok {
			return fmt.Errorf("invalid destination account balance format")
		}

		// Step 6: Calculate new balances
		newSourceBalance := new(big.Float).Sub(sourceBalance, amount)
		newDestBalance := new(big.Float).Add(destBalance, amount)

		// Step 7: Update source account balance
		// If this fails, transaction will rollback (no changes committed)
		if err := accountRepo.UpdateBalance(ctx, tx, req.SourceAccountID, newSourceBalance.Text('f', 5)); err != nil {
			return fmt.Errorf("failed to update source account balance: %w", err)
		}

		// Step 8: Update destination account balance
		// If this fails, transaction will rollback (source balance update also rolled back)
		if err := accountRepo.UpdateBalance(ctx, tx, req.DestinationAccountID, newDestBalance.Text('f', 5)); err != nil {
			return fmt.Errorf("failed to update destination account balance: %w", err)
		}

		// Step 9: Create transaction record
		// If this fails, transaction will rollback (both account updates also rolled back)
		transaction := &models.Transaction{
			SourceAccountID:      req.SourceAccountID,
			DestinationAccountID: req.DestinationAccountID,
			Amount:               req.Amount,
			Status:               "completed",
			CreatedAt:            time.Now(),
		}

		if err := transactionRepo.Create(ctx, transaction); err != nil {
			return fmt.Errorf("failed to create transaction record: %w", err)
		}

		// Return nil to commit the transaction
		// All steps succeeded, transaction will be committed atomically
		return nil
	})

	// If Transaction function returned an error, GORM automatically rolled back
	// Return the error to the caller
	return err
}
