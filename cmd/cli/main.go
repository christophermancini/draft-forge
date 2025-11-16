package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/yourusername/draft-forge/internal/db"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: cli <command> [args]")
		fmt.Println("\nAvailable commands:")
		fmt.Println("  migrate up    - Run all pending migrations")
		fmt.Println("  migrate down  - Rollback the last migration")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "migrate":
		if len(os.Args) < 3 {
			log.Fatal("Usage: cli migrate <up|down>")
		}
		handleMigrate(os.Args[2])
	default:
		log.Fatalf("Unknown command: %s", command)
	}
}

func handleMigrate(direction string) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	database, err := db.Connect(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := db.RunMigrations(database, direction); err != nil {
		log.Fatal(err)
	}
}
