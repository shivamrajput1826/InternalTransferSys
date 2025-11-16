package handlers

import (
	"internal-transfer-system/internal/models"
	"internal-transfer-system/internal/repository"
	"internal-transfer-system/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var transactionService = service.NewTransactionService()

// CreateTransfer handles POST /transactions
// Accepts JSON: { "source_account_id": 123, "destination_account_id": 456, "amount": "100.12345" }
// Processes the transaction to update account balances
// Returns: empty response (201) or error
func CreateTransfer(c *fiber.Ctx) error {
	var req models.CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate source_account_id
	if req.SourceAccountID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "source_account_id must be a positive integer",
		})
	}

	// Validate destination_account_id
	if req.DestinationAccountID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "destination_account_id must be a positive integer",
		})
	}

	// Validate accounts are different
	if req.SourceAccountID == req.DestinationAccountID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "source_account_id and destination_account_id cannot be the same",
		})
	}

	// Validate amount
	if req.Amount == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "amount is required",
		})
	}

	// Validate amount is a valid number
	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "amount must be a valid number",
		})
	}

	if amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "amount must be greater than zero",
		})
	}

	// Process transfer using service
	if err := transactionService.CreateTransfer(c.Context(), &req); err != nil {
		errMsg := err.Error()
		if err == repository.ErrInsufficientBalance {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "insufficient balance",
			})
		}
		// Check for account not found errors
		if err == repository.ErrAccountNotFound ||
			errMsg == "source account not found" ||
			errMsg == "destination account not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": errMsg,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errMsg,
		})
	}

	// Return empty response as per specification
	return c.Status(fiber.StatusCreated).Send(nil)
}
