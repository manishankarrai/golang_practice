package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"practice/api/models"
	"practice/api/services"

	"github.com/gin-gonic/gin"
)

func GetAvailableVehicles(c *gin.Context) {
	latitude, err := strconv.ParseFloat(c.Query("latitude"), 64)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "latitude is required"}); return }
	longitude, err := strconv.ParseFloat(c.Query("longitude"), 64)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "longitude is required"}); return }
	radius := 10.0
	if value := c.Query("radius_km"); value != "" { radius, err = strconv.ParseFloat(value, 64); if err != nil || radius <= 0 { c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "radius_km must be positive"}); return } }
	vehicles, err := services.GetAvailableVehicles(c.Request.Context(), latitude, longitude, radius)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"success": true, "data": vehicles})
}

func BookRide(c *gin.Context) {
	var req models.BookRideRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()}); return }
	ride, err := services.BookRide(c.Request.Context(), c.Param("customer_id"), req.VehicleID, req.Latitude, req.Longitude)
	if errors.Is(err, services.ErrCustomerNotFound) { c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "customer not found"}); return }
	if errors.Is(err, services.ErrVehicleUnavailable) { c.JSON(http.StatusConflict, gin.H{"success": false, "message": "vehicle is unavailable or outside the service area"}); return }
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "ride booked successfully", "data": ride})
}

func GetRides(c *gin.Context) {
	rides, err := services.GetRides(c.Request.Context())
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"success": true, "data": rides})
}

func GetRide(c *gin.Context) {
	ride, err := services.GetRideByID(c.Request.Context(), c.Param("id"))
	if errors.Is(err, services.ErrRideNotFound) { c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "ride not found"}); return }
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"success": true, "data": ride})
}
