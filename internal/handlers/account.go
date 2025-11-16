package handlers

import (
	"internal-transfer-system/internal/models"
	"internal-transfer-system/internal/repository"
	"internal-transfer-system/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var accountService = service.NewAccountService()

func CreateAccount(c *fiber.Ctx) error {
	var req models.CreateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation: account_id and initial_balance must be present
	if req.AccountID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "account_id must be a positive integer",
		})
	}

	if req.InitialBalance == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "initial_balance is required",
		})
	}

	// Create account using service
	if err := accountService.CreateAccount(c.Context(), &req); err != nil {
		if err == repository.ErrAccountAlreadyExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "account already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return empty response as per specification
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"account_id": req.AccountID,
		"balance":    req.InitialBalance,
		"message":    "Account created successfully",
	})
}
func GetAccount(c *fiber.Ctx) error {
	accountID, err := strconv.ParseInt(c.Params("account_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account_id",
		})
	}
	account, err := accountService.GetAccount(c.Context(), accountID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "account not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(account)
}
