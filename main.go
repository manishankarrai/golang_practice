package main

// Package main is the entry point for the golang_practice Gin-based web application.
// It loads environment configuration, initializes the shared MongoDB connection,
// sets up the router with middleware and routes, and starts the HTTP server.
import (
	"os"
	"practice/api/db"
	"practice/api/routes"
	"practice/config"

	"github.com/gin-gonic/gin"
)

// main bootstraps the application:
// - Loads .env based on APP_ENV
// - Connects to MongoDB via the db package
// - Creates a default Gin engine
// - Registers all API route groups (user, mobile, shop, location, etc.)
// - Starts listening on the configured port (default 8081)
func main() {
	config.LoadEnv()

	db.Connect() // initialize shared MongoDB client

	router := gin.Default()
	router.Static("/uploads", "./uploads")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})

	routes.RegisterRoutes(router) // initialized routes
	// read port from env or fallback to default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	router.Run(":" + port)
}
