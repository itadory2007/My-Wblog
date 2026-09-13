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
	postRepository := repository.NewPostRepository(conn)
	postAccessRepository := repository.NewPostAccessRepository(conn)
	commentRepository := repository.NewCommentRepository(conn)

	// Create services.
	userService := service.NewUserService(userRepository)
	sessionService := service.NewSessionService(sessionRepository)
	postService := service.NewPostService(postRepository,)
	postAccessService := service.NewPostAccessService(postAccessRepository, userRepository, postRepository,)
	commentService := service.NewCommentService(commentRepository)

	// Create handlers.
	authHandler := handlers.NewAuthHandler(userService, sessionService,)
	postHandler := handlers.NewPostHandler(postService, postAccessService, commentService,)

	// Create authentication middleware.
	authMiddleware := middleware.NewAuthMiddleware(sessionService,)

	// Create Echo server.
	e := echo.New()
	e.Renderer = handlers.NewTemplateRenderer()

	// Public routes.
	e.GET("/", postHandler.Feed, authMiddleware.RequireAuth,)
	
	e.POST("/signup", authHandler.Signup)
	e.GET("/signup", authHandler.ShowSignup)
	e.POST("/login", authHandler.Login)
	e.GET("/login", authHandler.ShowLogin)
	e.POST("/logout", authHandler.Logout)
	e.GET("/posts/new", postHandler.ShowCreatePost, authMiddleware.RequireAuth,)
	e.POST("/posts", postHandler.CreatePost, authMiddleware.RequireAuth,)
	e.GET("/weblog/:id", postHandler.GetPost,  authMiddleware.RequireAuth,)
	e.POST("/weblog/:id/access", postHandler.GrantAccess, authMiddleware.RequireAuth,)
	e.GET("/profile", func(c echo.Context) error {
			userID := c.Get(middleware.UserIDKey)
			return c.String(200, "Authenticated user ID: "+fmt.Sprint(userID),)
		},authMiddleware.RequireAuth,)

	// Start server.
	log.Println("Server started on http://localhost:8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}