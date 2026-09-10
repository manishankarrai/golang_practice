package main

import (
	"os"
	"practice/api/db"
	"practice/api/routes"
	"practice/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	db.Connect() // initialize shared MongoDB client

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})

	routes.RegisterRoutes(router) // initialized routes
	// read port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	router.Run()
}
