package service

import (
	"context"
	"errors"
	"fmt"
	"internal-transfer-system/internal/models"
	"internal-transfer-system/internal/repository"

	"gorm.io/gorm"
)

type AccountService struct {
	accountRepo *repository.AccountRepository
}

func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{
		accountRepo: repository.NewAccountRepository(db),
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, req *models.CreateAccountRequest) error {

	account := &models.Account{
		AccountID: req.AccountID,
		Balance:   req.InitialBalance,
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		if errors.Is(err, repository.ErrAccountAlreadyExists) {
			return repository.ErrAccountAlreadyExists
		}

		return fmt.Errorf("failed to create account: %w", err)
	}

	return nil
}

func (s *AccountService) GetAccount(ctx context.Context, accountID int64) (*models.AccountResponse, error) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, repository.ErrAccountNotFound) {
			return nil, repository.ErrAccountNotFound
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &models.AccountResponse{
		AccountID: account.AccountID,
		Balance:   account.Balance,
	}, nil
}
