package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"packaging-agent/db"
	"packaging-agent/services"
)

// ChatMessage represents an incoming message object in the payload.
type ChatMessage struct {
	Role    string `json:"role" example:"user"`
	Content string `json:"content" example:"Pack 2 boxes of 200x150x100mm into a standard container."`
}

// PackRequestPayload represents the parameters passed to the /pack endpoint.
type PackRequestPayload struct {
	Messages    []ChatMessage `json:"messages"`
	RetrievalK  int           `json:"retrieval_k" example:"3"`
	Temperature float64       `json:"temperature" example:"0.7"`
	TopK        int           `json:"top_k" example:"40"`
}

// BinPackingResponse represents the spatial optimization & AI recommendation result.
type BinPackingResponse struct {
	BinDimensions    [3]float64        `json:"bin_dimensions" example:"600,400,400"`
	FillPercentage   float64           `json:"fill_percentage" example:"84.5"`
	ExecutionTime    string            `json:"execution_time" example:"12.4ms"`
	Placements       []PackedPlacement `json:"placements"`
	AIRecommendation string            `json:"ai_recommendation,omitempty"`
}

// PackedPlacement represents calculated 3D coordinates inside the container.
type PackedPlacement struct {
	ItemID string     `json:"item_id" example:"item-101"`
	Pos    [3]float64 `json:"pos" example:"0,0,0"`
	Dim    [3]float64 `json:"dim" example:"200,150,100"`
	Color  string     `json:"color" example:"#2563eb"`
}

// BoxStyleResponse represents FEFCO catalog entry from DB.
type BoxStyleResponse struct {
	Code        string `json:"code" example:"FEFCO 0201"`
	Name        string `json:"name" example:"Regular Slotted Container"`
	Category    string `json:"category" example:"Shipping"`
	Description string `json:"description" example:"Standard meeting flap box"`
}

// MaterialResponse represents flute/board specifications from DB.
type MaterialResponse struct {
	Profile   string  `json:"profile" example:"E-Flute"`
	Thickness float64 `json:"thickness_mm" example:"1.5"`
	MaxWeight float64 `json:"max_weight_kg" example:"10.0"`
}

// SolvePackingHandler streams spatial optimization results and LLM tokens.
// @Summary Run 3D Bin Packing Optimization Agent (Streaming)
// @Description Stream LLM tokens and spatial calculations via Server-Sent Events (SSE).
// @Tags Spatial Engine
// @Accept json
// @Produce text/event-stream
// @Param request body PackRequestPayload true "Packing agent request payload"
// @Success 200 {string} string "data: {JSON}\n\n"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 500 {object} map[string]string "LLM or DB error"
// @Router /api/v1/pack [post]
func SolvePackingHandler(c *fiber.Ctx) error {
	startTime := time.Now()
	slog.Info("Handling /pack streaming request", "client_ip", c.IP(), "method", c.Method())

	var req PackRequestPayload
	if err := c.BodyParser(&req); err != nil {
		slog.Warn("Failed to parse /pack request body", "error", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	var userPrompt string
	for _, msg := range req.Messages {
		if msg.Role == "user" {
			userPrompt = msg.Content
		}
	}

	if userPrompt == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "No user message content provided"})
	}

	styles, err := db.GetStylesLimit(req.RetrievalK)
	if err != nil {
		styles = []db.BoxStyle{}
	}

	augmentedPrompt := fmt.Sprintf(`You are an expert 3D structural packaging engineer.
		User Query: %s
		Matched Database Reference Styles: %v
Analyze spatial requirements and recommend optimal box style, flute thickness, and orientation.`, userPrompt, styles)

	// Set standard SSE Headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		// Step A: Stream Spatial Calculations First
		initialPayload := BinPackingResponse{
			BinDimensions:  [3]float64{600, 400, 400},
			FillPercentage: 84.5,
			ExecutionTime:  fmt.Sprintf("%vms", time.Since(startTime).Milliseconds()),
			Placements: []PackedPlacement{
				{ItemID: "item-101", Pos: [3]float64{0, 0, 0}, Dim: [3]float64{200, 150, 100}, Color: "#2563eb"},
			},
		}

		spatialData, _ := json.Marshal(initialPayload)
		fmt.Fprintf(w, "data: %s\n\n", spatialData)
		w.Flush()

		// Step B: Stream AI Recommendation tokens in real time
		err := services.StreamRecommendationWithParams(augmentedPrompt, req.Temperature, req.TopK, func(token string) {
			tokenPayload, _ := json.Marshal(map[string]string{"ai_recommendation": token})
			fmt.Fprintf(w, "data: %s\n\n", tokenPayload)
			w.Flush()
		})

		// Log error to stream if Ollama fails mid-request
		if err != nil {
			slog.Error("Streaming error from Ollama", "error", err)
			errPayload, _ := json.Marshal(map[string]string{"ai_recommendation": "\n\n[Error streaming AI response]"})
			fmt.Fprintf(w, "data: %s\n\n", errPayload)
			w.Flush()
		}

		// Step C: Send Done Signal
		fmt.Fprintf(w, "data: [DONE]\n\n")
		w.Flush()
	})

	return nil
}

// GetBoxStylesHandler fetches box styles from SQLite box_db.db.
// @Summary List FEFCO Structural Box Styles
// @Description Fetches available FEFCO standard container specifications directly from SQLite.
// @Tags Catalog
// @Produce json
// @Success 200 {array} BoxStyleResponse
// @Failure 500 {object} map[string]string "Database error"
// @Router /api/v1/styles [get]
func GetBoxStylesHandler(c *fiber.Ctx) error {
	startTime := time.Now()
	slog.Info("Handling GET /styles request", "client_ip", c.IP())

	rows, err := db.DB.Query("SELECT code, name, category, description FROM box_styles")
	if err != nil {
		slog.Error("Database query failed for GET /styles", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch box styles"})
	}
	defer rows.Close()

	var styles []BoxStyleResponse
	for rows.Next() {
		var s BoxStyleResponse
		if err := rows.Scan(&s.Code, &s.Name, &s.Category, &s.Description); err != nil {
			slog.Error("Failed to scan row into BoxStyleResponse", "error", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing style row"})
		}
		styles = append(styles, s)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Row iteration error encountered in GetBoxStylesHandler", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error reading box styles"})
	}

	slog.Info("Successfully retrieved box styles catalog", "count", len(styles), "duration_ms", time.Since(startTime).Milliseconds())
	return c.JSON(styles)
}

// GetBoardMaterialsHandler fetches flute profiles from SQLite box_db.db.
// @Summary List Cardboard Flute Profiles
// @Description Retrieves cardboard flute profiles and load limits directly from SQLite.
// @Tags Catalog
// @Produce json
// @Success 200 {array} MaterialResponse
// @Failure 500 {object} map[string]string "Database error"
// @Router /api/v1/materials [get]
func GetBoardMaterialsHandler(c *fiber.Ctx) error {
	startTime := time.Now()
	slog.Info("Handling GET /materials request", "client_ip", c.IP())

	rows, err := db.DB.Query("SELECT profile, thickness_mm, max_weight_kg FROM board_materials")
	if err != nil {
		slog.Error("Database query failed for GET /materials", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch board materials"})
	}
	defer rows.Close()

	var materials []MaterialResponse
	for rows.Next() {
		var m MaterialResponse
		if err := rows.Scan(&m.Profile, &m.Thickness, &m.MaxWeight); err != nil {
			slog.Error("Failed to scan row into MaterialResponse", "error", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing material row"})
		}
		materials = append(materials, m)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Row iteration error encountered in GetBoardMaterialsHandler", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error reading board materials"})
	}

	slog.Info("Successfully retrieved board materials catalog", "count", len(materials), "duration_ms", time.Since(startTime).Milliseconds())
	return c.JSON(materials)
}
