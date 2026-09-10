package handlers

import (
	"net/http"

	"practice/api/services"

	"github.com/gin-gonic/gin"
)

// GetMobileCategories handles GET /api/mobile/categories.
func GetMobileCategories(c *gin.Context) {
	categories := services.GetMobileCategories()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "mobile categories fetched successfully",
		"data":    categories,
	})
}
