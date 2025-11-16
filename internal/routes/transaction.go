package routes

import (
	"internal-transfer-system/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupTransactionRoutes(app *fiber.App, prefix string) {
	transactionRoutes := app.Group(prefix + "/transactions")

	transactionRoutes.Post("/", handlers.CreateTransfer)
}
