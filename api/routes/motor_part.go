package routes

import (
	"practice/api/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterMotorParts registers motor-part CRUD endpoints.
func RegisterMotorParts(rg *gin.RouterGroup) {
	parts := rg.Group("/motor-parts")
	{
		parts.POST("", handlers.CreateMotorPart)
		parts.GET("", handlers.GetMotorParts)
		parts.GET("/:id", handlers.GetMotorPart)
		parts.PUT("/:id", handlers.UpdateMotorPart)
		parts.DELETE("/:id", handlers.DeleteMotorPart)
	}
}
