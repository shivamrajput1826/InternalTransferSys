package integration

import (
	"bytes"
	"encoding/json"
	"internal-transfer-system/config"
	"internal-transfer-system/internal/database"
	"internal-transfer-system/internal/models"
	"internal-transfer-system/internal/routes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	app *fiber.App
}

func (suite *APITestSuite) SetupSuite() {
	// Load test configuration
	config.LoadConfig()

	// Initialize database connection
	err := database.InitDB()
	suite.NoError(err, "Failed to initialize database")

	// Auto-migrate for tests
	err = database.AutoMigrate()
	suite.NoError(err, "Failed to auto-migrate database")

	// Create Fiber app
	suite.app = fiber.New(fiber.Config{
		BodyLimit:      2 * 1024 * 1024,
		Immutable:      true,
		ReadBufferSize: 4096,
	})

	// Setup routes
	API_PREFIX := "/api"
	routes.SetupRoutes(suite.app, API_PREFIX)
}

func (suite *APITestSuite) TearDownSuite() {
	// Clean up database connection
	database.CloseDB()
}

func (suite *APITestSuite) TearDownTest() {
	// Clean up test data after each test
	db := database.GetDB()
	if db != nil {
		db.Exec("TRUNCATE TABLE transactions CASCADE")
		db.Exec("TRUNCATE TABLE accounts CASCADE")
	}
}

// TestCreateAccount tests POST /api/accounts
func (suite *APITestSuite) TestCreateAccount() {
	tests := []struct {
		name           string
		payload        models.CreateAccountRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name: "successful account creation",
			payload: models.CreateAccountRequest{
				AccountID:      1,
				InitialBalance: "100.50",
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "duplicate account creation",
			payload: models.CreateAccountRequest{
				AccountID:      1,
				InitialBalance: "200.00",
			},
			expectedStatus: http.StatusConflict,
			expectError:    true,
		},
		{
			name: "invalid account_id",
			payload: models.CreateAccountRequest{
				AccountID:      0,
				InitialBalance: "100.50",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "missing initial_balance",
			payload: models.CreateAccountRequest{
				AccountID:      2,
				InitialBalance: "",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := suite.app.Test(req)
			suite.NoError(err)
			suite.Equal(tt.expectedStatus, resp.StatusCode)

			if !tt.expectError && tt.expectedStatus == http.StatusCreated {
				// Verify account was created by querying it
				getReq := httptest.NewRequest(http.MethodGet, "/api/accounts/1", nil)
				getResp, err := suite.app.Test(getReq)
				suite.NoError(err)
				suite.Equal(http.StatusOK, getResp.StatusCode)
			}
		})
	}
}

// TestGetAccount tests GET /api/accounts/{account_id}
func (suite *APITestSuite) TestGetAccount() {
	// First create an account
	createPayload := models.CreateAccountRequest{
		AccountID:      100,
		InitialBalance: "500.75",
	}
	body, _ := json.Marshal(createPayload)
	createReq := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBuffer(body))
	createReq.Header.Set("Content-Type", "application/json")
	createResp, _ := suite.app.Test(createReq)
	suite.Equal(http.StatusCreated, createResp.StatusCode)

	tests := []struct {
		name            string
		accountID       string
		expectedStatus  int
		expectedBalance string
	}{
		{
			name:            "successful account retrieval",
			accountID:       "100",
			expectedStatus:  http.StatusOK,
			expectedBalance: "500.75",
		},
		{
			name:           "account not found",
			accountID:      "999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid account_id format",
			accountID:      "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			req := httptest.NewRequest(http.MethodGet, "/api/accounts/"+tt.accountID, nil)
			resp, err := suite.app.Test(req)
			suite.NoError(err)
			suite.Equal(tt.expectedStatus, resp.StatusCode)

			if tt.expectedStatus == http.StatusOK {
				var account models.AccountResponse
				err := json.NewDecoder(resp.Body).Decode(&account)
				suite.NoError(err)
				suite.Equal(int64(100), account.AccountID)
				suite.Equal(tt.expectedBalance, account.Balance)
			}
		})
	}
}

// TestCreateTransfer tests POST /api/transactions
func (suite *APITestSuite) TestCreateTransfer() {
	// Setup: Create two accounts
	sourceAccount := models.CreateAccountRequest{
		AccountID:      200,
		InitialBalance: "1000.00",
	}
	destAccount := models.CreateAccountRequest{
		AccountID:      201,
		InitialBalance: "500.00",
	}

	body1, _ := json.Marshal(sourceAccount)
	req1 := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	suite.app.Test(req1)

	body2, _ := json.Marshal(destAccount)
	req2 := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	suite.app.Test(req2)

	tests := []struct {
		name           string
		payload        models.CreateTransactionRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name: "successful transfer",
			payload: models.CreateTransactionRequest{
				SourceAccountID:      200,
				DestinationAccountID: 201,
				Amount:               "100.50",
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "insufficient balance",
			payload: models.CreateTransactionRequest{
				SourceAccountID:      200,
				DestinationAccountID: 201,
				Amount:               "10000.00",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "source account not found",
			payload: models.CreateTransactionRequest{
				SourceAccountID:      999,
				DestinationAccountID: 201,
				Amount:               "100.00",
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name: "destination account not found",
			payload: models.CreateTransactionRequest{
				SourceAccountID:      200,
				DestinationAccountID: 998,
				Amount:               "100.00",
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name: "same source and destination",
			payload: models.CreateTransactionRequest{
				SourceAccountID:      200,
				DestinationAccountID: 200,
				Amount:               "100.00",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "invalid amount",
			payload: models.CreateTransactionRequest{
				SourceAccountID:      200,
				DestinationAccountID: 201,
				Amount:               "0",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := suite.app.Test(req)
			suite.NoError(err)
			suite.Equal(tt.expectedStatus, resp.StatusCode)

			// Verify balance changes for successful transfer
			if !tt.expectError && tt.expectedStatus == http.StatusCreated {
				// Check source account balance decreased
				getSourceReq := httptest.NewRequest(http.MethodGet, "/api/accounts/200", nil)
				getSourceResp, _ := suite.app.Test(getSourceReq)
				var sourceAccount models.AccountResponse
				json.NewDecoder(getSourceResp.Body).Decode(&sourceAccount)
				suite.NotEqual("1000.00", sourceAccount.Balance, "Source balance should have decreased")

				// Check destination account balance increased
				getDestReq := httptest.NewRequest(http.MethodGet, "/api/accounts/201", nil)
				getDestResp, _ := suite.app.Test(getDestReq)
				var destAccount models.AccountResponse
				json.NewDecoder(getDestResp.Body).Decode(&destAccount)
				suite.NotEqual("500.00", destAccount.Balance, "Destination balance should have increased")
			}
		})
	}
}

// TestTransferBalanceConsistency tests that balances are consistent after transfers
func (suite *APITestSuite) TestTransferBalanceConsistency() {
	// Create accounts
	sourceAccount := models.CreateAccountRequest{
		AccountID:      300,
		InitialBalance: "1000.00",
	}
	destAccount := models.CreateAccountRequest{
		AccountID:      301,
		InitialBalance: "500.00",
	}

	body1, _ := json.Marshal(sourceAccount)
	req1 := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	suite.app.Test(req1)

	body2, _ := json.Marshal(destAccount)
	req2 := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	suite.app.Test(req2)

	// Perform multiple transfers
	transfers := []models.CreateTransactionRequest{
		{SourceAccountID: 300, DestinationAccountID: 301, Amount: "100.00"},
		{SourceAccountID: 300, DestinationAccountID: 301, Amount: "50.00"},
		{SourceAccountID: 301, DestinationAccountID: 300, Amount: "25.00"},
	}

	for _, transfer := range transfers {
		body, _ := json.Marshal(transfer)
		req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := suite.app.Test(req)
		suite.Equal(http.StatusCreated, resp.StatusCode)
	}

	// Verify final balances
	getSourceReq := httptest.NewRequest(http.MethodGet, "/api/accounts/300", nil)
	getSourceResp, _ := suite.app.Test(getSourceReq)
	var sourceAccountResp models.AccountResponse
	json.NewDecoder(getSourceResp.Body).Decode(&sourceAccountResp)

	getDestReq := httptest.NewRequest(http.MethodGet, "/api/accounts/301", nil)
	getDestResp, _ := suite.app.Test(getDestReq)
	var destAccountResp models.AccountResponse
	json.NewDecoder(getDestResp.Body).Decode(&destAccountResp)

	// Total should remain constant: 1000 + 500 = 1500
	// After transfers: source should have 1000 - 100 - 50 + 25 = 875
	//                   dest should have 500 + 100 + 50 - 25 = 625
	// Total: 875 + 625 = 1500
	suite.Equal("875.00000", sourceAccountResp.Balance)
	suite.Equal("625.00000", destAccountResp.Balance)
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
