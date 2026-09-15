package routes

import (
	"practice/api/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterMobile registers all mobile related routes.
func RegisterMobile(rg *gin.RouterGroup) {
	mobile := rg.Group("/mobile")
	{
		// GET /mobile/categories - retrieves the list of supported mobile categories (e.g. Smartphone, Cellphone)
		mobile.GET("/categories", handlers.GetMobileCategories)

		smartphones := mobile.Group("/smartphones")
		{
			// POST /mobile/smartphones - creates a new smartphone record
			smartphones.POST("", handlers.CreateSmartphone)
			// GET /mobile/smartphones - lists all smartphones (optionally filtered by ?category=)
			smartphones.GET("", handlers.GetSmartphones)
			// GET /mobile/smartphones/:id - retrieves a single smartphone by its ID
			smartphones.GET("/:id", handlers.GetSmartphone)
			// PUT /mobile/smartphones/:id - updates an existing smartphone by its ID
			smartphones.PUT("/:id", handlers.UpdateSmartphone)
			// DELETE /mobile/smartphones/:id - deletes a smartphone by its ID
			smartphones.DELETE("/:id", handlers.DeleteSmartphone)
		}
	}
}
