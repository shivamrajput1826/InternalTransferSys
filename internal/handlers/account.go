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

var customLogger = logger.CreateLogger("Handler")

func CreateAccount(c *fiber.Ctx) error {
	log := customLogger.WithFiberContext(c)
	var req models.CreateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	log.LogRequest(logger.LogOptions{
		MethodName: "CreateAccount",
		Details:    req,
	})

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

	accountService := service.NewAccountService(database.GetDB())

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

	log.LogResponse(logger.LogOptions{
		MethodName: "CreateAccount",
		Details: map[string]interface{}{
			"account_id": req.AccountID,
			"balance":    req.InitialBalance,
			"message":    "Account created successfully",
		},
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"account_id": req.AccountID,
		"balance":    req.InitialBalance,
		"message":    "Account created successfully",
	})
}
func GetAccount(c *fiber.Ctx) error {
	log := customLogger.WithFiberContext(c)
	accountID, err := strconv.ParseInt(c.Params("account_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account_id",
		})
	}
	log.LogRequest(logger.LogOptions{
		MethodName: "GetAccount",
		Details:    accountID,
	})
	accountService := service.NewAccountService(database.GetDB())
	account, err := accountService.GetAccount(c.Context(), accountID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "account not found",
		})
	}
	log.LogResponse(logger.LogOptions{
		MethodName: "GetAccount",
		Details:    account,
	})
	return c.Status(fiber.StatusOK).JSON(account)
}
