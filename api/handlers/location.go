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
	// Create a variable to hold the location data sent in the request body.
	var req models.LocationRequest

	// Read the JSON request body and convert it into the LocationRequest struct.
	// ShouldBindJSON also validates the request according to the model's tags.
	if err := c.ShouldBindJSON(&req); err != nil {
		// Return HTTP 400 when the request body is invalid or cannot be parsed.
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		// Stop processing because there is no valid request to create.
		return
	}

	// Copy the values from the request struct into a Location model.
	// This is the object that will be passed to the service layer for creation.
	location := models.Location{
		// Set the location's name from the request.
		Name: req.Name,
		// Set the street or full address from the request.
		Address: req.Address,
		// Set the city from the request.
		City: req.City,
		// Set the state or province from the request.
		State: req.State,
		// Set the country from the request.
		Country: req.Country,
		// Set the postal or ZIP code from the request.
		ZipCode: req.ZipCode,
		// Set the geographic latitude from the request.
		Latitude: req.Latitude,
		// Set the geographic longitude from the request.
		Longitude: req.Longitude,
	}

	// Pass the request context and location model to the service layer.
	// The service layer performs the actual creation operation.
	created, err := services.CreateLocation(c.Request.Context(), location)
	if err != nil {
		// Return HTTP 500 when the service cannot create the location.
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		// Stop processing after sending the error response.
		return
	}

	// Return HTTP 201 to indicate that the location was created successfully.
	c.JSON(http.StatusCreated, gin.H{
		// Tell the client that the operation succeeded.
		"success": true,
		// Provide a human-readable success message.
		"message": "location created successfully",
		// Return the newly created location in the response data.
		"data": created,
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
