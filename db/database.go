package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type App struct {
	Queries *Queries
	Pool    *sql.DB
}

func InitDB() *App {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env")
	}

	dbURL := os.Getenv("TURSO_DATABASE_URL")
	dbToken := os.Getenv("TURSO_AUTH_TOKEN")
	connStr := fmt.Sprintf("%s?authToken=%s", dbURL, dbToken)
	if connStr == "" {
		log.Fatal("DB_URI environment variable is not set")
	}

	// Open connection
	db, err := sql.Open("libsql", connStr)
	if err != nil {
		log.Fatalf("Unable to open database: %v", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Println("⚡ Successfully connected to SQLite Cloud!")

	return &App{
		Queries: New(db),
		Pool:    db,
	}
}
