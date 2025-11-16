package repository

import (
	"context"
	"errors"
	"fmt"
	"internal-transfer-system/internal/models"

	"gorm.io/gorm"
)

var (
	ErrTransactionNotFound      = errors.New("transaction not found")
	ErrTransactionAlreadyExists = errors.New("transaction already exists")
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, transaction *models.Transaction) error {
	result := r.db.WithContext(ctx).Create(transaction)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrTransactionAlreadyExists
		}
		return fmt.Errorf("failed to create transaction: %w", result.Error)
	}
	return nil
}

func (r *TransactionRepository) GetByID(ctx context.Context, transactionID int64) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.WithContext(ctx).First(&transaction, transactionID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTransactionNotFound
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	return &transaction, nil
}
