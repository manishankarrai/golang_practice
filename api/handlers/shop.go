package handlers

import (
	"errors"
	"net/http"

	"practice/api/services"

	"github.com/gin-gonic/gin"
)

// GetShops handles GET /api/shops.
func GetShops(c *gin.Context) {
	shops, err := services.GetShops(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "shops fetched successfully",
		"data":    shops,
	})
}

// GetShop handles GET /api/shops/:id.
func GetShop(c *gin.Context) {
	shop, err := services.GetShopByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, services.ErrShopNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "shop not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "shop fetched successfully",
		"data":    shop,
	})
}
