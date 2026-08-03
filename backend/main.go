package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	// IMPORTANT: Import your auto-generated docs (will be created by `swag init`)
	_ "packaging-agent/docs"
)

// @title Packaging Agent Engine API
// @version 1.0
// @description AI-Powered Packaging Optimization Engine with FEFCO Dieline Generation.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	app := fiber.New()

	// Serve Swagger UI documentation endpoint
	app.Get("/swagger/*", fiberSwagger.HandlerDefault)

	// API Group
	v1 := app.Group("/api/v1")
	v1.Post("/recommend", HandleBoxRecommendation)

	log.Fatal(app.Listen(":8080"))
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
// @Router /recommend [post]
func HandleBoxRecommendation(c *fiber.Ctx) error {
	var input ProductInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Execution Logic (Postgres -> Gemini Temp 0.0 -> SVG Math)
	// ...
	return c.JSON(fiber.Map{"status": "success"})
}
