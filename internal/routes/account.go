package routes

import (
	"internal-transfer-system/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupAccountRoutes(app *fiber.App, prefix string) {
	accountRoutes := app.Group(prefix + "/accounts")

	accountRoutes.Post("/", handlers.CreateAccount)

	accountRoutes.Get("/:account_id", handlers.GetAccount)
}
