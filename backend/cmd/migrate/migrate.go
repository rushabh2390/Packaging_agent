package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"

	// Import your central config package
	"packaging-agent/config"
)

func main() {
	// 1. Load config
	cfg := config.LoadConfig()

	// 2. Open SQLite connection
	db, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		log.Fatalf("❌ Unable to open SQLite database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Verify connection
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("❌ Unable to connect to database: %v", err)
	}

	log.Printf("⚡ Connected to SQLite database at: %s\n", cfg.DBPath)

	// 3. Begin Transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("❌ Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// 4. Read SQL Files
	schemaSQL, err := os.ReadFile("db/schema.sql")
	if err != nil {
		log.Fatalf("❌ Failed to read db/schema.sql: %v", err)
	}

	seedSQL, err := os.ReadFile("db/seed.sql")
	if err != nil {
		log.Fatalf("❌ Failed to read db/seed.sql: %v", err)
	}

	// 5. Execute Schema & Seed
	log.Println("📦 Ensuring database tables exist...")
	if _, err := tx.ExecContext(ctx, string(schemaSQL)); err != nil {
		log.Fatalf("❌ Failed creating schema: %v", err)
	}

	log.Println("🌱 Seeding initial packaging data...")
	if _, err := tx.ExecContext(ctx, string(seedSQL)); err != nil {
		log.Fatalf("❌ Failed inserting seed data: %v", err)
	}

	// 6. Commit Transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("❌ Failed to commit transaction: %v", err)
	}

	fmt.Println("✅ SQLite database migration and seeding finished successfully!")
}
