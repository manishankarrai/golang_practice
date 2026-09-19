package handlers

import (
	"errors"
	"net/http"

	"practice/api/models"
	"practice/api/services"

	"github.com/gin-gonic/gin"
)

func CreateRider(c *gin.Context) {
	var req models.RiderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	rider, err := services.CreateRider(c.Request.Context(), models.Rider{Name: req.Name, Email: req.Email, Phone: req.Phone})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "rider created successfully", "data": rider})
}

func GetRiders(c *gin.Context) {
	riders, err := services.GetRiders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "riders fetched successfully", "data": riders})
}

func GetRider(c *gin.Context) {
	rider, err := services.GetRiderByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, services.ErrRiderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "rider not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "rider fetched successfully", "data": rider})
}

func UpdateRider(c *gin.Context) {
	var req models.RiderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	rider, err := services.UpdateRider(c.Request.Context(), c.Param("id"), models.Rider{Name: req.Name, Email: req.Email, Phone: req.Phone})
	if err != nil {
		if errors.Is(err, services.ErrRiderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "rider not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "rider updated successfully", "data": rider})
}

func DeleteRider(c *gin.Context) {
	if err := services.DeleteRider(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, services.ErrRiderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "rider not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "rider deleted successfully"})
}
