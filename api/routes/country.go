package routes

import (
	"practice/api/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterCountry registers all country related routes.
func RegisterCountry(rg *gin.RouterGroup) {
	countries := rg.Group("/countries")
	{
		countries.POST("", handlers.CreateCountry)
		countries.GET("", handlers.GetCountries)
		countries.GET("/:id", handlers.GetCountry)
		countries.PUT("/:id", handlers.UpdateCountry)
		countries.DELETE("/:id", handlers.DeleteCountry)
	}
}
