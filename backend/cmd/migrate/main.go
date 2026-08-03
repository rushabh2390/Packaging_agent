package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("❌ Error: DB_URL environment variable is not set")
	}

	ctx := context.Background()

	// 1. Establish connection pool
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("❌ Unable to connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("⚡ Connecting to PostgreSQL...")

	// 2. Begin Transaction
	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx) // Rollback if execution fails before Commit()

	// 3. Read SQL Files
	schemaSQL, err := os.ReadFile("db/schema.sql")
	if err != nil {
		log.Fatalf("❌ Failed to read db/schema.sql: %v", err)
	}

	seedSQL, err := os.ReadFile("db/seed.sql")
	if err != nil {
		log.Fatalf("❌ Failed to read db/seed.sql: %v", err)
	}

	// 4. Execute Table Creation (IF NOT EXISTS)
	log.Println("📦 Ensuring database tables exist...")
	if _, err := tx.Exec(ctx, string(schemaSQL)); err != nil {
		log.Fatalf("❌ Failed creating schema: %v", err)
	}

	// 5. Execute Data Seeding (ON CONFLICT DO NOTHING)
	log.Println("🌱 Seeding initial packaging data...")
	if _, err := tx.Exec(ctx, string(seedSQL)); err != nil {
		log.Fatalf("❌ Failed inserting seed data: %v", err)
	}

	// 6. Commit Transaction
	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("❌ Failed to commit transaction: %v", err)
	}

	fmt.Println("✅ Database migration and seeding finished successfully!")
}
