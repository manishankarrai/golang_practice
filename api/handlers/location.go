package handlers

import (
	"errors"
	"net/http"

	"practice/api/models"
	"practice/api/services"

	"github.com/gin-gonic/gin"
)

// CreateLocation handles POST /api/locations.
func CreateLocation(c *gin.Context) {
	var req models.LocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	location := models.Location{
		Name:      req.Name,
		Address:   req.Address,
		City:      req.City,
		State:     req.State,
		Country:   req.Country,
		ZipCode:   req.ZipCode,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	created, err := services.CreateLocation(c.Request.Context(), location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "location created successfully",
		"data":    created,
	})
}

// GetLocations handles GET /api/locations with an optional ?city= filter.
func GetLocations(c *gin.Context) {
	city := c.Query("city")

	locations, err := services.GetLocations(c.Request.Context(), city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "locations fetched successfully",
		"data":    locations,
	})
}

// GetLocation handles GET /api/locations/:id.
func GetLocation(c *gin.Context) {
	location, err := services.GetLocationByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, services.ErrLocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "location not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "location fetched successfully",
		"data":    location,
	})
}

// UpdateLocation handles PUT /api/locations/:id.
func UpdateLocation(c *gin.Context) {
	var req models.LocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	location := models.Location{
		Name:      req.Name,
		Address:   req.Address,
		City:      req.City,
		State:     req.State,
		Country:   req.Country,
		ZipCode:   req.ZipCode,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	updated, err := services.UpdateLocation(c.Request.Context(), c.Param("id"), location)
	if err != nil {
		if errors.Is(err, services.ErrLocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "location not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "location updated successfully",
		"data":    updated,
	})
}

// DeleteLocation handles DELETE /api/locations/:id.
func DeleteLocation(c *gin.Context) {
	if err := services.DeleteLocation(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, services.ErrLocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "location not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "location deleted successfully"})
}
