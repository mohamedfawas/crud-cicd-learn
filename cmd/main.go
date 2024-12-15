package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/crud-cicd-learn/internal/handler"
	"github.com/mohamedfawas/crud-cicd-learn/internal/service"
)

func main() {
	// Initialize our service
	userService := service.NewUserService()

	// Initialize handler
	userHandler := handler.NewUserHandler(userService)

	// Set up Gin router
	router := gin.Default()

	// Define routes
	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)
	router.GET("/users", userHandler.GetAllUsers)

	// Start server
	log.Println("Server starting on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
