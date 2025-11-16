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

type ServiceTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *ServiceTestSuite) SetupSuite() {
	// Use in-memory SQLite database for unit tests
	var err error
	suite.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.NoError(err)

	// Auto-migrate
	suite.db.AutoMigrate(&models.Account{}, &models.Transaction{})

	// Initialize services with test database
	// Note: We need to modify services to accept DB, or use dependency injection
	// For now, we'll test the logic directly
}

func (suite *ServiceTestSuite) TearDownTest() {
	// Clean up after each test
	suite.db.Exec("DELETE FROM transactions")
	suite.db.Exec("DELETE FROM accounts")
}

func (suite *ServiceTestSuite) TestTransactionService_CreateTransfer() {
	accountRepo := repository.NewAccountRepository(suite.db)
	ctx := context.Background()

	// Setup: Create accounts
	sourceAccount := &models.Account{
		AccountID: 200,
		Balance:   "1000",
	}
	destAccount := &models.Account{
		AccountID: 201,
		Balance:   "500",
	}
	accountRepo.Create(ctx, sourceAccount)
	accountRepo.Create(ctx, destAccount)

	// Test successful transfer using transaction
	err := suite.db.Transaction(func(tx *gorm.DB) error {
		txAccountRepo := repository.NewAccountRepository(tx)

		// Get balances with lock (verify they exist)
		_, err := txAccountRepo.GetBalanceForUpdate(ctx, tx, 200)
		suite.NoError(err)

		_, err = txAccountRepo.GetBalanceForUpdate(ctx, tx, 201)
		suite.NoError(err)

		// Update balances
		err = txAccountRepo.UpdateBalance(ctx, tx, 200, "900")
		suite.NoError(err)

		err = txAccountRepo.UpdateBalance(ctx, tx, 201, "600")
		suite.NoError(err)

		// Create transaction record
		transaction := &models.Transaction{
			SourceAccountID:      200,
			DestinationAccountID: 201,
			Amount:               "100",
			Status:               "completed",
		}
		txTransactionRepo := repository.NewTransactionRepository(tx)
		return txTransactionRepo.Create(ctx, transaction)
	})

	suite.NoError(err)

	// Verify balances were updated
	source, _ := accountRepo.GetByID(ctx, 200)
	dest, _ := accountRepo.GetByID(ctx, 201)
	suite.Equal("900", source.Balance)
	suite.Equal("600", dest.Balance)
}

// TestTransactionService_TransferAtomicity tests that transfers are atomic
func (suite *ServiceTestSuite) TestTransactionService_TransferAtomicity() {
	accountRepo := repository.NewAccountRepository(suite.db)
	ctx := context.Background()

	// Create account with insufficient balance
	sourceAccount := &models.Account{
		AccountID: 300,
		Balance:   "100",
	}
	destAccount := &models.Account{
		AccountID: 301,
		Balance:   "500",
	}
	accountRepo.Create(ctx, sourceAccount)
	accountRepo.Create(ctx, destAccount)

	// Try to transfer more than available (should fail and rollback)
	err := suite.db.Transaction(func(tx *gorm.DB) error {
		txAccountRepo := repository.NewAccountRepository(tx)

		// Get balance
		balance, err := txAccountRepo.GetBalanceForUpdate(ctx, tx, 300)
		suite.NoError(err)
		suite.Equal("100", balance)

		// Try to update with insufficient balance - simulate error
		// In real service, this would check balance first
		return repository.ErrInsufficientBalance
	})

	suite.Error(err)
	suite.Equal(repository.ErrInsufficientBalance, err)

	// Verify balances were NOT changed (transaction rolled back)
	source, _ := accountRepo.GetByID(ctx, 300)
	dest, _ := accountRepo.GetByID(ctx, 301)
	suite.Equal("100", source.Balance, "Source balance should not have changed")
	suite.Equal("500", dest.Balance, "Destination balance should not have changed")
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

// Helper function for simple assertions
func TestAccountService_Simple(t *testing.T) {
	// This is a simple test without suite
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	db.AutoMigrate(&models.Account{})

	accountRepo := repository.NewAccountRepository(db)
	ctx := context.Background()

	// Test account creation
	account := &models.Account{
		AccountID: 1,
		Balance:   "100",
	}
	err = accountRepo.Create(ctx, account)
	assert.NoError(t, err)

	// Test account retrieval
	retrieved, err := accountRepo.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), retrieved.AccountID)
	assert.Equal(t, "100", retrieved.Balance)
}
