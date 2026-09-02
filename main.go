package main

import (
	"context"
	"log"
	"os"

	"weblog/internal/database"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// Load environment variables.
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env")
	}

	// Get database URL.
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Connect to PostgreSQL.
	conn, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer conn.Close(context.Background())

	// Create Echo server.
	e := echo.New()

	// Test route.
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Weblog is running!")
	})

	// Start server.
	log.Println("Server started on http://localhost:8080")

	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
