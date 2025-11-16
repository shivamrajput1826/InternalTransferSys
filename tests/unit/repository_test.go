package unit

import (
	"context"
	"internal-transfer-system/internal/models"
	"internal-transfer-system/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type RepositoryTestSuite struct {
	suite.Suite
	db              *gorm.DB
	accountRepo     *repository.AccountRepository
	transactionRepo *repository.TransactionRepository
}

func (suite *RepositoryTestSuite) SetupSuite() {
	var err error
	suite.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.NoError(err)

	// Auto-migrate
	suite.db.AutoMigrate(&models.Account{}, &models.Transaction{})

	// Initialize repositories
	suite.accountRepo = repository.NewAccountRepository(suite.db)
	suite.transactionRepo = repository.NewTransactionRepository(suite.db)
}

func (suite *RepositoryTestSuite) TearDownTest() {
	// Clean up after each test
	suite.db.Exec("DELETE FROM transactions")
	suite.db.Exec("DELETE FROM accounts")
}

// TestAccountRepository_Create tests account creation
func (suite *RepositoryTestSuite) TestAccountRepository_Create() {
	ctx := context.Background()

	tests := []struct {
		name        string
		account     *models.Account
		expectError bool
	}{
		{
			name: "successful account creation",
			account: &models.Account{
				AccountID: 1,
				Balance:   "100.50",
			},
			expectError: false,
		},
		{
			name: "duplicate account",
			account: &models.Account{
				AccountID: 1,
				Balance:   "200.00",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := suite.accountRepo.Create(ctx, tt.account)

			if tt.expectError {
				suite.Error(err)
				suite.Equal(repository.ErrAccountAlreadyExists, err)
			} else {
				suite.NoError(err)
			}
		})
	}
}

// TestAccountRepository_GetByID tests account retrieval
func (suite *RepositoryTestSuite) TestAccountRepository_GetByID() {
	ctx := context.Background()

	// Create an account first
	account := &models.Account{
		AccountID: 100,
		Balance:   "500.75",
	}
	err := suite.accountRepo.Create(ctx, account)
	suite.NoError(err)

	tests := []struct {
		name        string
		accountID   int64
		expectError bool
		errorType   error
	}{
		{
			name:        "successful retrieval",
			accountID:   100,
			expectError: false,
		},
		{
			name:        "account not found",
			accountID:   999,
			expectError: true,
			errorType:   repository.ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			account, err := suite.accountRepo.GetByID(ctx, tt.accountID)

			if tt.expectError {
				suite.Error(err)
				suite.Nil(account)
				if tt.errorType != nil {
					suite.Equal(tt.errorType, err)
				}
			} else {
				suite.NoError(err)
				suite.NotNil(account)
				suite.Equal(tt.accountID, account.AccountID)
				suite.Equal("500.75", account.Balance)
			}
		})
	}
}

// TestAccountRepository_Exists tests account existence check
func (suite *RepositoryTestSuite) TestAccountRepository_Exists() {
	ctx := context.Background()

	// Create an account
	account := &models.Account{
		AccountID: 200,
		Balance:   "300.00",
	}
	suite.accountRepo.Create(ctx, account)

	tests := []struct {
		name        string
		accountID   int64
		expected    bool
		expectError bool
	}{
		{
			name:        "account exists",
			accountID:   200,
			expected:    true,
			expectError: false,
		},
		{
			name:        "account does not exist",
			accountID:   999,
			expected:    false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			exists, err := suite.accountRepo.Exists(ctx, tt.accountID)

			if tt.expectError {
				suite.Error(err)
			} else {
				suite.NoError(err)
				suite.Equal(tt.expected, exists)
			}
		})
	}
}

// TestAccountRepository_UpdateBalance tests balance update
func (suite *RepositoryTestSuite) TestAccountRepository_UpdateBalance() {
	ctx := context.Background()

	// Create an account
	account := &models.Account{
		AccountID: 300,
		Balance:   "1000.00",
	}
	suite.accountRepo.Create(ctx, account)

	tests := []struct {
		name        string
		accountID   int64
		newBalance  string
		expectError bool
		errorType   error
	}{
		{
			name:        "successful update",
			accountID:   300,
			newBalance:  "1500.00",
			expectError: false,
		},
		{
			name:        "account not found",
			accountID:   999,
			newBalance:  "500.00",
			expectError: true,
			errorType:   repository.ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := suite.accountRepo.UpdateBalance(ctx, nil, tt.accountID, tt.newBalance)

			if tt.expectError {
				suite.Error(err)
				if tt.errorType != nil {
					suite.Equal(tt.errorType, err)
				}
			} else {
				suite.NoError(err)
				// Verify balance was updated
				updatedAccount, _ := suite.accountRepo.GetByID(ctx, tt.accountID)
				suite.Equal(tt.newBalance, updatedAccount.Balance)
			}
		})
	}
}

// TestAccountRepository_GetBalanceForUpdate tests row-level locking
func (suite *RepositoryTestSuite) TestAccountRepository_GetBalanceForUpdate() {
	ctx := context.Background()

	// Create an account
	account := &models.Account{
		AccountID: 400,
		Balance:   "750.25",
	}
	suite.accountRepo.Create(ctx, account)

	// Start a transaction
	tx := suite.db.Begin()
	defer tx.Rollback()

	accountRepo := repository.NewAccountRepository(tx)
	balance, err := accountRepo.GetBalanceForUpdate(ctx, tx, 400)

	suite.NoError(err)
	suite.Equal("750.25", balance)
}

// TestTransactionRepository_Create tests transaction creation
func (suite *RepositoryTestSuite) TestTransactionRepository_Create() {
	ctx := context.Background()

	tests := []struct {
		name        string
		transaction *models.Transaction
		expectError bool
	}{
		{
			name: "successful transaction creation",
			transaction: &models.Transaction{
				SourceAccountID:      1,
				DestinationAccountID: 2,
				Amount:               "100.50",
				Status:               "completed",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := suite.transactionRepo.Create(ctx, tt.transaction)

			if tt.expectError {
				suite.Error(err)
			} else {
				suite.NoError(err)
				suite.Greater(tt.transaction.ID, int64(0), "Transaction ID should be set")
			}
		})
	}
}

// TestTransactionRepository_GetByID tests transaction retrieval
func (suite *RepositoryTestSuite) TestTransactionRepository_GetByID() {
	ctx := context.Background()

	// Create a transaction
	transaction := &models.Transaction{
		SourceAccountID:      1,
		DestinationAccountID: 2,
		Amount:               "200.00",
		Status:               "completed",
	}
	suite.transactionRepo.Create(ctx, transaction)

	tests := []struct {
		name          string
		transactionID int64
		expectError   bool
		errorType     error
	}{
		{
			name:          "successful retrieval",
			transactionID: transaction.ID,
			expectError:   false,
		},
		{
			name:          "transaction not found",
			transactionID: 999,
			expectError:   true,
			errorType:     repository.ErrTransactionNotFound,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			transaction, err := suite.transactionRepo.GetByID(ctx, tt.transactionID)

			if tt.expectError {
				suite.Error(err)
				suite.Nil(transaction)
				if tt.errorType != nil {
					suite.Equal(tt.errorType, err)
				}
			} else {
				suite.NoError(err)
				suite.NotNil(transaction)
				suite.Equal(tt.transactionID, transaction.ID)
			}
		})
	}
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

// Simple test without suite
func TestAccountRepository_Simple(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	db.AutoMigrate(&models.Account{})

	repo := repository.NewAccountRepository(db)
	ctx := context.Background()

	account := &models.Account{
		AccountID: 1,
		Balance:   "100.00",
	}

	err = repo.Create(ctx, account)
	assert.NoError(t, err)

	retrieved, err := repo.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), retrieved.AccountID)
	assert.Equal(t, "100.00", retrieved.Balance)
}
