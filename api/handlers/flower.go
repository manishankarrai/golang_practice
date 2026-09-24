package handlers

import (
	"errors"
	"net/http"

	"practice/api/models"
	"practice/api/services"

	"github.com/gin-gonic/gin"
)

// CreateFlower handles POST /api/flowers.
func CreateFlower(c *gin.Context) {
	var req models.FlowerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	flower := models.Flower{
		Name:     req.Name,
		Type:     req.Type,
		Color:    req.Color,
		Location: req.Location,
		Pricing:  req.Pricing,
		Details:  req.Details,
	}

	created, err := services.CreateFlower(c.Request.Context(), flower)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "flower created successfully",
		"data":    created,
	})
}

// GetFlowers handles GET /api/flowers.
func GetFlowers(c *gin.Context) {
	flowers, err := services.GetFlowers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "flowers fetched successfully",
		"data":    flowers,
	})
}

// GetFlower handles GET /api/flowers/:id.
func GetFlower(c *gin.Context) {
	flower, err := services.GetFlowerByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, services.ErrFlowerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "flower not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "flower fetched successfully", "data": flower})
}

// UpdateFlower handles PUT /api/flowers/:id.
func UpdateFlower(c *gin.Context) {
	var req models.FlowerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	flower := models.Flower{
		Name:     req.Name,
		Type:     req.Type,
		Color:    req.Color,
		Location: req.Location,
		Pricing:  req.Pricing,
		Details:  req.Details,
	}

	updated, err := services.UpdateFlower(c.Request.Context(), c.Param("id"), flower)
	if err != nil {
		if errors.Is(err, services.ErrFlowerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "flower not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "flower updated successfully", "data": updated})
}

// DeleteFlower handles DELETE /api/flowers/:id.
func DeleteFlower(c *gin.Context) {
	if err := services.DeleteFlower(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, services.ErrFlowerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "flower not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "flower deleted successfully"})
}
