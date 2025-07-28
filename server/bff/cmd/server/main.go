package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/tadasy/todo-app/server/bff/internal/api/handlers"
	customMiddleware "github.com/tadasy/todo-app/server/bff/internal/api/middleware"
	"github.com/tadasy/todo-app/server/bff/internal/clients"
)

func main() {
	// Initialize gRPC clients
	userClient, err := clients.NewUserServiceClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	defer userClient.Close()

	todoClient, err := clients.NewTodoServiceClient("localhost:50052")
	if err != nil {
		log.Fatalf("Failed to connect to todo service: %v", err)
	}
	defer todoClient.Close()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userClient)
	todoHandler := handlers.NewTodoHandler(todoClient)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	// Public routes
	e.POST("/api/auth/register", authHandler.Register)
	e.POST("/api/auth/login", authHandler.Login)

	// Protected routes
	api := e.Group("/api")
	api.Use(customMiddleware.JWTMiddleware)

	// Todo routes
	api.POST("/todos", todoHandler.CreateTodo)
	api.GET("/todos", todoHandler.ListTodos)
	api.GET("/todos/:id", todoHandler.GetTodo)
	api.PUT("/todos/:id", todoHandler.UpdateTodo)
	api.PUT("/todos/:id/complete", todoHandler.MarkTodoComplete)
	api.DELETE("/todos/:id", todoHandler.DeleteTodo)

	// Start server
	log.Println("BFF server starting on port 8080...")
	log.Fatal(e.Start(":8080"))
}
