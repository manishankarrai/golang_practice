package handlers

import (
	"errors"
	"log"
	"net/http"

	"practice/api/models"
	"practice/api/services"

	"github.com/gin-gonic/gin"
)

func CreateMotorPart(c *gin.Context) {
	var req models.MotorPartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	part, err := services.CreateMotorPart(c.Request.Context(), models.MotorPart{Name: req.Name, Description: req.Description, Price: req.Price})
	if err != nil {
		log.Printf("CreateMotorPart: service error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "motor part created successfully", "data": part})
}

func GetMotorParts(c *gin.Context) {
	parts, err := services.GetMotorParts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "motor parts fetched successfully", "data": parts})
}

func GetMotorPart(c *gin.Context) {
	part, err := services.GetMotorPartByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, services.ErrMotorPartNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "motor part not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "motor part fetched successfully", "data": part})
}

func UpdateMotorPart(c *gin.Context) {
	var req models.MotorPartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	part, err := services.UpdateMotorPart(c.Request.Context(), c.Param("id"), models.MotorPart{Name: req.Name, Description: req.Description, Price: req.Price})
	if err != nil {
		if errors.Is(err, services.ErrMotorPartNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "motor part not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "motor part updated successfully", "data": part})
}

func DeleteMotorPart(c *gin.Context) {
	if err := services.DeleteMotorPart(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, services.ErrMotorPartNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "motor part not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "motor part deleted successfully"})
}
