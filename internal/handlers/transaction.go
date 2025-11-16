package handlers

import (
	"internal-transfer-system/internal/database"
	"internal-transfer-system/internal/models"
	"internal-transfer-system/internal/repository"
	"internal-transfer-system/internal/service"
	"internal-transfer-system/logger"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateTransfer(c *fiber.Ctx) error {
	log := customLogger.WithFiberContext(c)
	var req models.CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	log.LogRequest(logger.LogOptions{
		MethodName: "CreateTransfer",
		Details:    req,
	})

	if req.SourceAccountID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "source_account_id must be a positive integer",
		})
	}

	if req.DestinationAccountID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "destination_account_id must be a positive integer",
		})
	}

	if req.SourceAccountID == req.DestinationAccountID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "source_account_id and destination_account_id cannot be the same",
		})
	}

	if req.Amount == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "amount is required",
		})
	}

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

	transactionService := service.NewTransactionService(database.GetDB())
	if err := transactionService.CreateTransfer(c.Context(), &req); err != nil {
		errMsg := err.Error()
		if err == repository.ErrInsufficientBalance {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "insufficient balance",
			})
		}
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
	log.LogResponse(logger.LogOptions{
		MethodName: "CreateTransfer",
		Details:    "Successfull Transaction",
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "successful transaction",
	})
}
