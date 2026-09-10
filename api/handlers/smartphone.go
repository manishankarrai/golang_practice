package handlers

import (
	"errors"
	"net/http"

	"practice/api/models"
	"practice/api/services"

	"github.com/gin-gonic/gin"
)

// CreateSmartphone handles POST /api/mobile/smartphones.
func CreateSmartphone(c *gin.Context) {
	var req models.SmartphoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	phone := models.Smartphone{
		Name:     req.Name,
		Brand:    req.Brand,
		Category: req.Category,
		Price:    req.Price,
		Stock:    req.Stock,
	}

	created, err := services.CreateSmartphone(c.Request.Context(), phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "smartphone created successfully",
		"data":    created,
	})
}

// GetSmartphones handles GET /api/mobile/smartphones with an optional ?category= filter.
func GetSmartphones(c *gin.Context) {
	category := c.Query("category")

	phones, err := services.GetSmartphones(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "smartphones fetched successfully",
		"data":    phones,
	})
}

// GetSmartphone handles GET /api/mobile/smartphones/:id.
func GetSmartphone(c *gin.Context) {
	phone, err := services.GetSmartphoneByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, services.ErrSmartphoneNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "smartphone not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "smartphone fetched successfully",
		"data":    phone,
	})
}

// UpdateSmartphone handles PUT /api/mobile/smartphones/:id.
func UpdateSmartphone(c *gin.Context) {
	var req models.SmartphoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	phone := models.Smartphone{
		Name:     req.Name,
		Brand:    req.Brand,
		Category: req.Category,
		Price:    req.Price,
		Stock:    req.Stock,
	}

	updated, err := services.UpdateSmartphone(c.Request.Context(), c.Param("id"), phone)
	if err != nil {
		if errors.Is(err, services.ErrSmartphoneNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "smartphone not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "smartphone updated successfully",
		"data":    updated,
	})
}

// DeleteSmartphone handles DELETE /api/mobile/smartphones/:id.
func DeleteSmartphone(c *gin.Context) {
	if err := services.DeleteSmartphone(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, services.ErrSmartphoneNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "smartphone not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "smartphone deleted successfully"})
}
