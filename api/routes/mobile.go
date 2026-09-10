package routes

import (
	"practice/api/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterMobile registers all mobile related routes.
func RegisterMobile(rg *gin.RouterGroup) {
	mobile := rg.Group("/mobile")
	{
		mobile.GET("/categories", handlers.GetMobileCategories)

		smartphones := mobile.Group("/smartphones")
		{
			smartphones.POST("", handlers.CreateSmartphone)
			smartphones.GET("", handlers.GetSmartphones)
			smartphones.GET("/:id", handlers.GetSmartphone)
			smartphones.PUT("/:id", handlers.UpdateSmartphone)
			smartphones.DELETE("/:id", handlers.DeleteSmartphone)
		}
	}
}
