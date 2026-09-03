package main

import (
	"context"
	"log"
	"os"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"weblog/internal/database"
	"weblog/internal/handlers"
	"weblog/internal/repositories"
	"weblog/internal/services"
	"weblog/internal/middleware"
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

	// Create repositories.
	userRepository := repository.NewUserRepository(conn)
	sessionRepository := repository.NewSessionRepository(conn)

	// Create services.
	userService := service.NewUserService(userRepository)
	sessionService := service.NewSessionService(sessionRepository)

	// Create handlers.
	authHandler := handlers.NewAuthHandler(userService, sessionService,)

	// Create authentication middleware.
	authMiddleware := middleware.NewAuthMiddleware(sessionService,)

	// Create Echo server.
	e := echo.New()

	// Public routes.
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Weblog is running!")
	})

	e.POST("/signup", authHandler.Signup)
	e.GET("/signup", authHandler.ShowSignup)
	e.POST("/login", authHandler.Login)
	e.GET("/login", authHandler.ShowLogin)
	e.POST("/logout", authHandler.Logout)

	// Protected test route.
	e.GET("/profile", func(c echo.Context) error {
			userID := c.Get(middleware.UserIDKey)
			return c.String(200, "Authenticated user ID: "+fmt.Sprint(userID),)
		},
		authMiddleware.RequireAuth,
	)

	// Start server.
	log.Println("Server started on http://localhost:8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}