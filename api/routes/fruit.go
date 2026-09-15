package routes

import (
	"practice/api/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterFruit registers all fruit CRUD routes.
func RegisterFruit(rg *gin.RouterGroup) {
	fruits := rg.Group("/fruits")
	{
		fruits.POST("", handlers.CreateFruit)
		fruits.GET("", handlers.GetFruits)
		fruits.GET("/:id", handlers.GetFruit)
		fruits.PUT("/:id", handlers.UpdateFruit)
		fruits.DELETE("/:id", handlers.DeleteFruit)
	}
}
