package main

import (
	"os"
	"time"

	_ "packaging-agent/docs"
	"packaging-agent/handlers"
	"packaging-agent/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

// ProductInput defines the schema for incoming packaging specifications
type ProductInput struct {
	Width  float64 `json:"width" example:"200"`
	Height float64 `json:"height" example:"150"`
	Length float64 `json:"length" example:"300"`
	Weight float64 `json:"weight" example:"2.5"`
}

// PackagingRecommendation defines the response schema for the endpoint
type PackagingRecommendation struct {
	Status     string `json:"status" example:"success"`
	BoxStyle   string `json:"box_style,omitempty" example:"FEFCO 0201"`
	DielineSVG string `json:"dieline_svg,omitempty"`
}

func main() {
	// 1. Initialize Structured Logger
	logger.InitLogger()
	logger.Log.Info("Starting Packaging Agent Backend Service...")

	// 2. Initialize Fiber App
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Get("/swagger/*", swagger.HandlerDefault)
	// 2. Attach CORS Middleware (Allows Next.js / Frontend requests)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Or "http://localhost:3000" for strict production control
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
	}))

	// 3. Attach Request Logging Middleware
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "[HTTP] ${time} | ${status} | ${latency} | ${ip} | ${method} ${path}\n",
	}))
	// API Grouping
	v1 := app.Group("/api/v1")
	v1.Post("/pack", handlers.SolvePackingHandler)
	v1.Post("/recommend", HandleBoxRecommendation)

	// 4. Start Server with Graceful Error Logging
	serverAddr := ":8080"
	logger.Log.Info("Server listener starting", "address", serverAddr)

	if err := app.Listen(serverAddr); err != nil {
		logger.Log.Error("Backend server shut down unexpectedly", "error", err)
		os.Exit(1)
	}
}

// HandleBoxRecommendation handles product packaging evaluation
// @Summary Calculate optimal packaging box and dieline
// @Description Accepts product specs, checks Postgres for clearance rules, and queries Gemini @ Temp 0.0
// @Tags Packaging
// @Accept json
// @Produce json
// @Param request body ProductInput true "Product Specification Payload"
// @Success 200 {object} PackagingRecommendation "Recommended FEFCO box and SVG layout"
// @Failure 400 {object} map[string]string "Invalid input error"
// @Failure 500 {object} map[string]string "Internal processing error"
// @Router /api/v1/recommend [post]
func HandleBoxRecommendation(c *fiber.Ctx) error {
	startTime := time.Now()

	logger.Log.Info("Processing /api/v1/recommend request",
		"client_ip", c.IP(),
		"method", c.Method(),
	)

	var input ProductInput
	if err := c.BodyParser(&input); err != nil {
		logger.Log.Warn("Failed to parse recommendation request body",
			"error", err,
			"client_ip", c.IP(),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	logger.Log.Debug("Parsed product input specifications",
		"length", input.Length,
		"width", input.Width,
		"height", input.Height,
		"weight", input.Weight,
	)

	// Execution Logic (Postgres -> Gemini Temp 0.0 -> SVG Math)
	// Example placeholder step:
	// recommend, err := services.CalculateBoxRecommendation(input)
	// if err != nil {
	// 	logger.Log.Error("Error calculating box recommendation", "error", err)
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to calculate recommendation"})
	// }

	logger.Log.Info("Successfully generated box recommendation",
		"duration_ms", time.Since(startTime).Milliseconds(),
		"client_ip", c.IP(),
	)

	return c.JSON(fiber.Map{
		"status": "success",
	})
}
