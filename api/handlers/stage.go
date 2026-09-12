package handlers

import (
	"errors"
	"net/http"

	"practice/api/models"
	"practice/api/services"

	"github.com/gin-gonic/gin"
)

// CreateStage handles POST /api/stages.
func CreateStage(c *gin.Context) {
	var req models.StageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	stage := models.Stage{
		Name:        req.Name,
		Description: req.Description,
		Sequence:    req.Sequence,
		Status:      req.Status,
	}

	created, err := services.CreateStage(c.Request.Context(), stage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "stage created successfully",
		"data":    created,
	})
}

// GetStages handles GET /api/stages.
func GetStages(c *gin.Context) {
	stages, err := services.GetStages(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "stages fetched successfully",
		"data":    stages,
	})
}

// GetStage handles GET /api/stages/:id.
func GetStage(c *gin.Context) {
	stage, err := services.GetStageByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, services.ErrStageNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "stage not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "stage fetched successfully",
		"data":    stage,
	})
}

// UpdateStage handles PUT /api/stages/:id.
func UpdateStage(c *gin.Context) {
	var req models.StageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	stage := models.Stage{
		Name:        req.Name,
		Description: req.Description,
		Sequence:    req.Sequence,
		Status:      req.Status,
	}

	updated, err := services.UpdateStage(c.Request.Context(), c.Param("id"), stage)
	if err != nil {
		if errors.Is(err, services.ErrStageNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "stage not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "stage updated successfully",
		"data":    updated,
	})
}

// DeleteStage handles DELETE /api/stages/:id.
func DeleteStage(c *gin.Context) {
	if err := services.DeleteStage(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, services.ErrStageNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "stage not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "stage deleted successfully"})
}
