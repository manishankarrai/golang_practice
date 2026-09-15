package handlers

import (
	"errors"
	"net/http"

	"practice/api/models"
	"practice/api/services"

	"github.com/gin-gonic/gin"
)

// CreateFruit handles POST /api/fruits.
func CreateFruit(c *gin.Context) {
	var req models.FruitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	fruit := models.Fruit{
		Name:     req.Name,
		Type:     req.Type,
		Location: req.Location,
		Pricing:  req.Pricing,
		Details:  req.Details,
	}

	created, err := services.CreateFruit(c.Request.Context(), fruit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "fruit created successfully",
		"data":    created,
	})
}

// GetFruits handles GET /api/fruits.
func GetFruits(c *gin.Context) {
	fruits, err := services.GetFruits(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "fruits fetched successfully",
		"data":    fruits,
	})
}

// GetFruit handles GET /api/fruits/:id.
func GetFruit(c *gin.Context) {
	fruit, err := services.GetFruitByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, services.ErrFruitNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "fruit not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "fruit fetched successfully",
		"data":    fruit,
	})
}

// UpdateFruit handles PUT /api/fruits/:id.
func UpdateFruit(c *gin.Context) {
	var req models.FruitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	fruit := models.Fruit{
		Name:     req.Name,
		Type:     req.Type,
		Location: req.Location,
		Pricing:  req.Pricing,
		Details:  req.Details,
	}

	updated, err := services.UpdateFruit(c.Request.Context(), c.Param("id"), fruit)
	if err != nil {
		if errors.Is(err, services.ErrFruitNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "fruit not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "fruit updated successfully",
		"data":    updated,
	})
}

// DeleteFruit handles DELETE /api/fruits/:id.
func DeleteFruit(c *gin.Context) {
	if err := services.DeleteFruit(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, services.ErrFruitNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "fruit not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "fruit deleted successfully"})
}
