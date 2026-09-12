package handlers

import (
	"errors"
	"net/http"

	"practice/api/models"
	"practice/api/services"

	"github.com/gin-gonic/gin"
)

// CreateCountry handles POST /api/countries.
func CreateCountry(c *gin.Context) {
	var req models.CountryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	country := models.Country{
		Name: req.Name,
		Code: req.Code,
	}

	created, err := services.CreateCountry(c.Request.Context(), country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "country created successfully",
		"data":    created,
	})
}

// GetCountries handles GET /api/countries.
func GetCountries(c *gin.Context) {
	countries, err := services.GetCountries(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "countries fetched successfully",
		"data":    countries,
	})
}

// GetCountry handles GET /api/countries/:id.
func GetCountry(c *gin.Context) {
	country, err := services.GetCountryByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, services.ErrCountryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "country not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "country fetched successfully",
		"data":    country,
	})
}

// UpdateCountry handles PUT /api/countries/:id.
func UpdateCountry(c *gin.Context) {
	var req models.CountryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	country := models.Country{
		Name: req.Name,
		Code: req.Code,
	}

	updated, err := services.UpdateCountry(c.Request.Context(), c.Param("id"), country)
	if err != nil {
		if errors.Is(err, services.ErrCountryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "country not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "country updated successfully",
		"data":    updated,
	})
}

// DeleteCountry handles DELETE /api/countries/:id.
func DeleteCountry(c *gin.Context) {
	if err := services.DeleteCountry(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, services.ErrCountryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "country not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "country deleted successfully"})
}
