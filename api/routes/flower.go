package routes

import (
	"practice/api/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterFlower registers all flower CRUD routes.
func RegisterFlower(rg *gin.RouterGroup) {
	flowers := rg.Group("/flowers")
	{
		flowers.POST("", handlers.CreateFlower)
		flowers.GET("", handlers.GetFlowers)
		flowers.GET("/:id", handlers.GetFlower)
		flowers.PUT("/:id", handlers.UpdateFlower)
		flowers.DELETE("/:id", handlers.DeleteFlower)
	}
}
