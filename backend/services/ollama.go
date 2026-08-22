package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"packaging-agent/config"
	"time"
)

// Ensure System field is present in request struct
type OllamaParamRequest struct {
	Model   string        `json:"model"`
	System  string        `json:"system"`
	Prompt  string        `json:"prompt"`
	Stream  bool          `json:"stream"`
	Options OllamaOptions `json:"options"`
}

type OllamaOptions struct {
	Temperature float64 `json:"temperature"`
	TopK        int     `json:"top_k"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

const SystemPrompt = `You are a 3D Packaging & Box Optimization Engineering Agent.

When recommending box sizes, FEFCO styles, or material specs, you MUST strictly structure your output in GitHub Flavored Markdown (GFM):

1. **Summary Recommendation**: Direct 1-2 sentence core verdict.
2. **Dimension Comparison Table**: A clean Markdown table comparing Product Input Size (mm) vs Minimum Buffer (mm) vs Outer Box Dimensions (L×W×H mm).
3. **FEFCO & Material Details**: Use bold titles and itemized bullet points for Flute Thickness, FEFCO Style Code, Load Capacity, and Box Features.

Example Table Format:
| Parameter | Product Size | Recommended Outer Box |
| :--- | :--- | :--- |
| **Length** | 300 mm | **320 mm** |
| **Width** | 200 mm | **220 mm** |
| **Height** | 120 mm | **150 mm** |

Do NOT return continuous unformatted prose or walls of plain text.`

// StreamRecommendationWithParams streams individual LLM tokens from local Ollama to a callback function.
func StreamRecommendationWithParams(prompt string, temp float64, topK int, onToken func(string)) error {
	startTime := time.Now()
	modelName := "qwen2.5-coder:3b"

	slog.Info("Initiating Ollama streaming recommendation",
		"model", modelName,
		"temperature", temp,
		"top_k", topK,
		"prompt_len", len(prompt),
	)

	cfg := config.LoadConfig()
	reqBody := OllamaParamRequest{
		Model:  modelName,
		System: SystemPrompt,
		Prompt: prompt,
		Stream: true, // Enabled token streaming
		Options: OllamaOptions{
			Temperature: temp,
			TopK:        topK,
		},
	}

	targetURL := fmt.Sprintf("%s/api/generate", cfg.OllamaURL)
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		slog.Error("Failed to marshal request body for Ollama", "error", err, "model", modelName)
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Post(targetURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("Network error reaching local Ollama instance",
			"error", err,
			"target_url", targetURL,
			"duration_ms", time.Since(startTime).Milliseconds(),
		)
		return fmt.Errorf("failed to reach local Ollama instance: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		slog.Error("Ollama API returned non-200 HTTP status code",
			"status_code", resp.StatusCode,
			"response_body", string(bodyBytes),
			"duration_ms", time.Since(startTime).Milliseconds(),
		)
		return fmt.Errorf("ollama API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	// Read line-delimited JSON objects from Ollama stream
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var chunk OllamaResponse
		if err := json.Unmarshal(line, &chunk); err != nil {
			slog.Warn("Failed to decode token chunk from Ollama", "error", err, "line", string(line))
			continue
		}

		// Stream individual token back via callback
		if chunk.Response != "" {
			onToken(chunk.Response)
		}
	}

	if err := scanner.Err(); err != nil {
		slog.Error("Error encountered while reading Ollama stream", "error", err)
		return fmt.Errorf("error reading stream: %w", err)
	}

	slog.Info("Successfully completed Ollama token stream", "duration_ms", time.Since(startTime).Milliseconds())
	return nil
}
