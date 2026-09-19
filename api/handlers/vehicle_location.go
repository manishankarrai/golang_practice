package handlers

import (
	"errors"
	"net/http"

	"practice/api/models"
	"practice/api/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func UpdateVehicleLocation(c *gin.Context) {
	var req models.VehicleLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	riderID, err := bson.ObjectIDFromHex(req.RiderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid rider_id"})
		return
	}
	if _, err := services.GetRiderByID(c.Request.Context(), req.RiderID); err != nil {
		if errors.Is(err, services.ErrRiderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "rider not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	location, err := services.SaveVehicleLocation(c.Request.Context(), c.Param("id"), models.VehicleLocation{RiderID: &riderID, Latitude: req.Latitude, Longitude: req.Longitude})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "vehicle location updated successfully", "data": location})
}

func GetVehicleLocation(c *gin.Context) {
	location, fromCache, err := services.GetVehicleLocation(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, services.ErrVehicleLocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "vehicle location not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "vehicle location fetched successfully", "cache_hit": fromCache, "data": location})
}
