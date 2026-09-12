package handlers

import (
	"errors"
	"net/http"

	"practice/api/models"
	"practice/api/services"

	"github.com/gin-gonic/gin"
)

// CreateState handles POST /api/states.
func CreateState(c *gin.Context) {
	var req models.StateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	state := models.State{
		Name:    req.Name,
		Code:    req.Code,
		Country: req.Country,
	}

	created, err := services.CreateState(c.Request.Context(), state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "state created successfully",
		"data":    created,
	})
}

// GetStates handles GET /api/states with an optional ?country= filter.
func GetStates(c *gin.Context) {
	country := c.Query("country")

	states, err := services.GetStates(c.Request.Context(), country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "states fetched successfully",
		"data":    states,
	})
}

// GetState handles GET /api/states/:id.
func GetState(c *gin.Context) {
	state, err := services.GetStateByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, services.ErrStateNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "state not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "state fetched successfully",
		"data":    state,
	})
}

// UpdateState handles PUT /api/states/:id.
func UpdateState(c *gin.Context) {
	var req models.StateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	state := models.State{
		Name:    req.Name,
		Code:    req.Code,
		Country: req.Country,
	}

	updated, err := services.UpdateState(c.Request.Context(), c.Param("id"), state)
	if err != nil {
		if errors.Is(err, services.ErrStateNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "state not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "state updated successfully",
		"data":    updated,
	})
}

// DeleteState handles DELETE /api/states/:id.
func DeleteState(c *gin.Context) {
	if err := services.DeleteState(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, services.ErrStateNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "state not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "state deleted successfully"})
}
