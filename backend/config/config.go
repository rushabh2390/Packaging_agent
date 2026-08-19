package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPath    string
	OllamaURL string
}

func LoadConfig() *Config {
	if err := godotenv.Load(".env", "../.env"); err != nil {
		log.Println("ℹ️  No local .env file loaded. Relying on system/container environment variables.")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = os.Getenv("DB_URL")
	}
	if dbPath == "" {
		dbPath = "./box_db.db" // Fallback default
	}

	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434" // Default local Ollama endpoint
	}

	return &Config{
		DBPath:    dbPath,
		OllamaURL: ollamaURL,
	}
}
