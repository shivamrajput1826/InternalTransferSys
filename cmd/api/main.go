package main

import (
	"internal-transfer-system/config"
	"internal-transfer-system/internal/database"
	"internal-transfer-system/internal/routes"
	"internal-transfer-system/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

var customLogger = logger.CreateLogger("internal-system")

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database connection
	if err := database.InitDB(); err != nil {
		customLogger.Error(logger.LogOptions{
			MethodName: "main",
			Details:    err,
		})
		panic(err)
	}
	defer func() {
		if err := database.CloseDB(); err != nil {
			customLogger.Error(logger.LogOptions{
				MethodName: "main",
				Details:    err,
			})
		}
	}()

	// Auto-migrate database schema (creates/updates tables based on models)
	// Note: For production, use SQL migrations instead
	if err := database.AutoMigrate(); err != nil {
		customLogger.Error(logger.LogOptions{
			MethodName: "main",
			Details:    err,
		})
		panic(err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit:      2 * 1024 * 1024,
		Immutable:      true,
		ReadBufferSize: 4096,
	})

	// Setup routes
	API_PREFIX := "/api"
	routes.SetupRoutes(app, API_PREFIX)

	// Graceful shutdown
	go func() {
		port := config.GetPort()
		if err := app.Listen(":" + port); err != nil {
			customLogger.Error(logger.LogOptions{
				MethodName: "main",
				Details:    err,
			})
			panic(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	customLogger.Debug("Shutting down server...", logger.LogOptions{
		MethodName: "main",
		Details:    "Shutting down server...",
	})

	if err := app.Shutdown(); err != nil {
		customLogger.Error(logger.LogOptions{
			MethodName: "main",
			Details:    err,
		})
	}
}
