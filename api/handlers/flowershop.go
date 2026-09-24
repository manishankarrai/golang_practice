package handlers

import (
	"errors"
	"net/http"

	"practice/api/models"
	"practice/api/services"

	"github.com/gin-gonic/gin"
)

// CreateFlowerShop handles POST /api/flowershops.
func CreateFlowerShop(c *gin.Context) {
	var req models.FlowerShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	shop := models.FlowerShop{
		Name:     req.Name,
		Owner:    req.Owner,
		Location: req.Location,
		Phone:    req.Phone,
		Flowers:  make([]models.Flower, 0),
	}

	created, err := services.CreateFlowerShop(c.Request.Context(), shop)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "flower shop created successfully",
		"data":    created,
	})
}

// GetFlowerShops handles GET /api/flowershops.
func GetFlowerShops(c *gin.Context) {
	shops, err := services.GetFlowerShops(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "flower shops fetched successfully",
		"data":    shops,
	})
}

// GetFlowerShop handles GET /api/flowershops/:id.
func GetFlowerShop(c *gin.Context) {
	shop, err := services.GetFlowerShopByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, services.ErrFlowerShopNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "flower shop not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "flower shop fetched successfully", "data": shop})
}

// UpdateFlowerShop handles PUT /api/flowershops/:id.
func UpdateFlowerShop(c *gin.Context) {
	var req models.FlowerShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	shop := models.FlowerShop{
		Name:     req.Name,
		Owner:    req.Owner,
		Location: req.Location,
		Phone:    req.Phone,
	}

	updated, err := services.UpdateFlowerShop(c.Request.Context(), c.Param("id"), shop)
	if err != nil {
		if errors.Is(err, services.ErrFlowerShopNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "flower shop not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "flower shop updated successfully", "data": updated})
}

// DeleteFlowerShop handles DELETE /api/flowershops/:id.
func DeleteFlowerShop(c *gin.Context) {
	if err := services.DeleteFlowerShop(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, services.ErrFlowerShopNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "flower shop not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "flower shop deleted successfully"})
}

// AddFlowerToShop handles POST /api/flowershops/:id/flowers.
func AddFlowerToShop(c *gin.Context) {
	var req models.FlowerShopFlowerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	updated, err := services.AddFlowerToShop(c.Request.Context(), c.Param("id"), req.FlowerID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrFlowerShopNotFound):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "flower shop not found"})
		case errors.Is(err, services.ErrFlowerNotFound):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "flower not found"})
		case errors.Is(err, services.ErrFlowerAlreadyInShop):
			c.JSON(http.StatusConflict, gin.H{"success": false, "message": "flower already exists in shop"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "flower added to shop successfully", "data": updated})
}

// RemoveFlowerFromShop handles DELETE /api/flowershops/:id/flowers/:flowerid.
func RemoveFlowerFromShop(c *gin.Context) {
	updated, err := services.RemoveFlowerFromShop(c.Request.Context(), c.Param("id"), c.Param("flowerid"))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrFlowerShopNotFound):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "flower shop not found"})
		case errors.Is(err, services.ErrFlowerNotInShop):
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "flower not found in shop"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "flower removed from shop successfully", "data": updated})
}
