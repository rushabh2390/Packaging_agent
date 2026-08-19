package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

// InitLogger initializes a structured JSON logger for production
func InitLogger() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug, // Adjust level based on ENV (Info/Debug/Error)
		AddSource: true,            // Include file line numbers for errors
	})

	Log = slog.New(handler)
	slog.SetDefault(Log)
}
