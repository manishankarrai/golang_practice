package handlers

import (
	"errors"
	"net/http"

	"practice/api/models"
	"practice/api/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func registrationFromRequest(req models.RegistrationRequest) (models.Registration, error) {
	vehicleID, err := bson.ObjectIDFromHex(req.VehicleID)
	if err != nil { return models.Registration{}, err }
	return models.Registration{VehicleID: vehicleID, OwnerName: req.OwnerName, OwnerEmail: req.OwnerEmail, OwnerPhone: req.OwnerPhone, RegistrationNumber: req.RegistrationNumber, RegistrationDate: req.RegistrationDate, ExpiryDate: req.ExpiryDate, OwnerAddress: req.OwnerAddress}, nil
}

func CreateRegistration(c *gin.Context) {
	var req models.RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()}); return }
	registration, err := registrationFromRequest(req)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "vehicle_id must be a valid object ID"}); return }
	created, err := services.CreateRegistration(c.Request.Context(), registration)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "registration created successfully", "data": created})
}

func GetRegistrations(c *gin.Context) {
	registrations, err := services.GetRegistrations(c.Request.Context())
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "registrations fetched successfully", "data": registrations})
}

func GetRegistration(c *gin.Context) {
	registration, err := services.GetRegistrationByID(c.Request.Context(), c.Param("id"))
	if err != nil { if errors.Is(err, services.ErrRegistrationNotFound) { c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "registration not found"}); return }; c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "registration fetched successfully", "data": registration})
}

func UpdateRegistration(c *gin.Context) {
	var req models.RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()}); return }
	registration, err := registrationFromRequest(req)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "vehicle_id must be a valid object ID"}); return }
	updated, err := services.UpdateRegistration(c.Request.Context(), c.Param("id"), registration)
	if err != nil { if errors.Is(err, services.ErrRegistrationNotFound) { c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "registration not found"}); return }; c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "registration updated successfully", "data": updated})
}

func DeleteRegistration(c *gin.Context) {
	if err := services.DeleteRegistration(c.Request.Context(), c.Param("id")); err != nil { if errors.Is(err, services.ErrRegistrationNotFound) { c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "registration not found"}); return }; c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "registration deleted successfully"})
}
