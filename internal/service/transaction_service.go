package service

import (
	"context"
	"errors"
	"fmt"
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

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{
		accountRepo:     repository.NewAccountRepository(db),
		transactionRepo: repository.NewTransactionRepository(db),
		db:              db,
	}
}

func (s *TransactionService) CreateTransfer(ctx context.Context, req *models.CreateTransactionRequest) error {
	amount, ok := new(big.Float).SetString(req.Amount)
	if !ok {
		return fmt.Errorf("invalid amount format")
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		accountRepo := repository.NewAccountRepository(tx)
		transactionRepo := repository.NewTransactionRepository(tx)

		sourceBalanceStr, err := accountRepo.GetBalanceForUpdate(ctx, tx, req.SourceAccountID)
		if err != nil {
			if errors.Is(err, repository.ErrAccountNotFound) {
				return fmt.Errorf("source account not found")
			}
			return fmt.Errorf("failed to get source account balance: %w", err)
		}

		sourceBalance, ok := new(big.Float).SetString(sourceBalanceStr)
		if !ok {
			return fmt.Errorf("invalid source account balance format")
		}

		if sourceBalance.Cmp(amount) < 0 {
			return repository.ErrInsufficientBalance
		}

		destBalanceStr, err := accountRepo.GetBalanceForUpdate(ctx, tx, req.DestinationAccountID)
		if err != nil {
			if errors.Is(err, repository.ErrAccountNotFound) {
				return fmt.Errorf("destination account not found")
			}
			return fmt.Errorf("failed to get destination account balance: %w", err)
		}

		destBalance, ok := new(big.Float).SetString(destBalanceStr)
		if !ok {
			return fmt.Errorf("invalid destination account balance format")
		}

		newSourceBalance := new(big.Float).Sub(sourceBalance, amount)
		newDestBalance := new(big.Float).Add(destBalance, amount)

		if err := accountRepo.UpdateBalance(ctx, tx, req.SourceAccountID, newSourceBalance.Text('f', 5)); err != nil {
			return fmt.Errorf("failed to update source account balance: %w", err)
		}

		if err := accountRepo.UpdateBalance(ctx, tx, req.DestinationAccountID, newDestBalance.Text('f', 5)); err != nil {
			return fmt.Errorf("failed to update destination account balance: %w", err)
		}

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

		return nil
	})

	return err
}
