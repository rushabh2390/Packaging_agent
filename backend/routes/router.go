package routes

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"packaging-agent/handlers"
)

func SetupRoutes(app *fiber.App) {
	// Add global HTTP access logging middleware for all incoming requests
	app.Use(logger.New(logger.Config{
		Format: "[HTTP] ${time} | ${status} | ${latency} | ${ip} | ${method} ${path}\n",
	}))

	slog.Info("Registering application routes...")

	// Top-level endpoint matching your cURL directly: http://localhost:8080/pack
	app.Post("/pack", handlers.SolvePackingHandler)

	// API v1 grouping
	api := app.Group("/api/v1")

	// Routes under /api/v1/
	api.Post("/pack", handlers.SolvePackingHandler)
	api.Get("/styles", handlers.GetBoxStylesHandler)
	api.Get("/materials", handlers.GetBoardMaterialsHandler)

	slog.Info("Routes registered successfully")
}
