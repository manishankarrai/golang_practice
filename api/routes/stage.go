package routes

import (
	"practice/api/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterStage registers all stage related routes.
func RegisterStage(rg *gin.RouterGroup) {
	stages := rg.Group("/stages")
	{
		stages.POST("", handlers.CreateStage)
		stages.GET("", handlers.GetStages)
		stages.GET("/:id", handlers.GetStage)
		stages.PUT("/:id", handlers.UpdateStage)
		stages.DELETE("/:id", handlers.DeleteStage)
	}
}
