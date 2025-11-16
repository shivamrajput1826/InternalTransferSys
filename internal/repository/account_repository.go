package repository

import (
	"context"
	"errors"
	"fmt"
	"internal-transfer-system/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrAccountNotFound      = errors.New("account not found")
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrInsufficientBalance  = errors.New("insufficient balance")
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, account *models.Account) error {
	result := r.db.WithContext(ctx).Create(account)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrAccountAlreadyExists
		}
		return fmt.Errorf("failed to create account: %w", result.Error)
	}
	return nil
}

func (r *AccountRepository) GetByID(ctx context.Context, accountID int64) (*models.Account, error) {
	var account models.Account
	err := r.db.WithContext(ctx).First(&account, accountID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return &account, nil
}

func (r *AccountRepository) Exists(ctx context.Context, accountID int64) (bool, error) {
	var account models.Account
	err := r.db.WithContext(ctx).Select("1").First(&account, accountID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check account existence: %w", err)
	}
	return true, nil
}

func (r *AccountRepository) UpdateBalance(ctx context.Context, tx *gorm.DB, accountID int64, newBalance string) error {
	var db = r.db
	if tx != nil {
		db = tx
	}
	result := db.WithContext(ctx).Model(&models.Account{}).Where("account_id = ?", accountID).Update("balance", newBalance)
	if result.Error != nil {
		return fmt.Errorf("failed to update balance: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrAccountNotFound
	}
	return nil
}

func (r *AccountRepository) GetBalanceForUpdate(ctx context.Context, tx *gorm.DB, accountID int64) (string, error) {
	var account models.Account
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&account, accountID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrAccountNotFound
		}
		return "", fmt.Errorf("failed to get balance for update: %w", err)
	}
	return account.Balance, nil
}
