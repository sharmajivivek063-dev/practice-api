package main

import (
	"fmt"

	"practice-api/connection"
	"practice-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	// Connect database
	connection.ConnectDB()
	defer connection.CloseDB()

	// Initialize Gin
	router := gin.Default()

	// Register routes
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.TaskRoutes(router)

	// Start server
	fmt.Println("Server running on port 8080")

	router.Run(":8080")
}
