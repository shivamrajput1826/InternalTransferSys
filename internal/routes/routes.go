package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, apiPrefix string) {
	SetupAccountRoutes(app, apiPrefix)

	SetupTransactionRoutes(app, apiPrefix)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})
}
