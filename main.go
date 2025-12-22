package main

import (
	"practice/api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})

	routes.RegisterRoutes(router) // initialized routes

	router.Run()
}
