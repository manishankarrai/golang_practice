package routes

import (
	"practice/api/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterLocation registers all location related routes.
func RegisterLocation(rg *gin.RouterGroup) {
	locations := rg.Group("/locations")
	{
		locations.POST("", handlers.CreateLocation)
		locations.GET("", handlers.GetLocations)
		locations.GET("/:id", handlers.GetLocation)
		locations.PUT("/:id", handlers.UpdateLocation)
		locations.DELETE("/:id", handlers.DeleteLocation)
	}
}
