package handlers

import (
	"errors"
	"log"
	"net/http"

	"practice/api/models"
	"practice/api/services"

	"github.com/gin-gonic/gin"
)

func CreateMotorVehicle(c *gin.Context) {
	log.Printf("CreateMotorVehicle: request received")
	var req models.MotorVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("CreateMotorVehicle: invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	vehicle, err := services.CreateMotorVehicle(c.Request.Context(), models.MotorVehicle{Make: req.Make, Model: req.Model, Year: req.Year, Color: req.Color, VehicleType: req.VehicleType, VIN: req.VIN, EngineNo: req.EngineNo})
	if err != nil {
		log.Printf("CreateMotorVehicle: service error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	log.Printf("CreateMotorVehicle: vehicle created: %s", vehicle.ID.Hex())
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "motor vehicle created successfully", "data": vehicle})
}

func GetMotorVehicles(c *gin.Context) {
	log.Printf("GetMotorVehicles: request received")
	vehicles, err := services.GetMotorVehicles(c.Request.Context())
	if err != nil {
		log.Printf("GetMotorVehicles: service error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	log.Printf("GetMotorVehicles: returned %d vehicles", len(vehicles))
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "motor vehicles fetched successfully", "data": vehicles})
}

func GetMotorVehicle(c *gin.Context) {
	id := c.Param("id")
	log.Printf("GetMotorVehicle: request received for id=%s", id)
	vehicle, err := services.GetMotorVehicleByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("GetMotorVehicle: service error for id=%s: %v", id, err)
		if errors.Is(err, services.ErrMotorVehicleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "motor vehicle not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	log.Printf("GetMotorVehicle: vehicle found: %s", vehicle.ID.Hex())
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "motor vehicle fetched successfully", "data": vehicle})
}

func UpdateMotorVehicle(c *gin.Context) {
	id := c.Param("id")
	log.Printf("UpdateMotorVehicle: request received for id=%s", id)
	var req models.MotorVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("UpdateMotorVehicle: invalid request for id=%s: %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	vehicle, err := services.UpdateMotorVehicle(c.Request.Context(), id, models.MotorVehicle{Make: req.Make, Model: req.Model, Year: req.Year, Color: req.Color, VehicleType: req.VehicleType, VIN: req.VIN, EngineNo: req.EngineNo})
	if err != nil {
		log.Printf("UpdateMotorVehicle: service error for id=%s: %v", id, err)
		if errors.Is(err, services.ErrMotorVehicleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "motor vehicle not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	log.Printf("UpdateMotorVehicle: vehicle updated: %s", vehicle.ID.Hex())
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "motor vehicle updated successfully", "data": vehicle})
}

func AttachRiderToMotorVehicle(c *gin.Context) {
	id := c.Param("id")
	log.Printf("AttachRiderToMotorVehicle: request received for vehicle id=%s", id)
	var req models.AttachRiderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("AttachRiderToMotorVehicle: invalid request for vehicle id=%s: %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	if _, err := services.GetRiderByID(c.Request.Context(), req.RiderID); err != nil {
		log.Printf("AttachRiderToMotorVehicle: rider lookup error for rider id=%s: %v", req.RiderID, err)
		if errors.Is(err, services.ErrRiderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "rider not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	vehicle, err := services.SetVehicleRider(c.Request.Context(), id, req.RiderID)
	if err != nil {
		log.Printf("AttachRiderToMotorVehicle: service error for vehicle id=%s: %v", id, err)
		if errors.Is(err, services.ErrMotorVehicleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "motor vehicle not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	log.Printf("AttachRiderToMotorVehicle: rider %s attached to vehicle %s", req.RiderID, vehicle.ID.Hex())
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "rider attached to motor vehicle successfully", "data": vehicle})
}

func DetachRiderFromMotorVehicle(c *gin.Context) {
	id := c.Param("id")
	log.Printf("DetachRiderFromMotorVehicle: request received for vehicle id=%s", id)
	vehicle, err := services.ClearVehicleRider(c.Request.Context(), id)
	if err != nil {
		log.Printf("DetachRiderFromMotorVehicle: service error for vehicle id=%s: %v", id, err)
		if errors.Is(err, services.ErrMotorVehicleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "motor vehicle not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	log.Printf("DetachRiderFromMotorVehicle: rider detached from vehicle %s", vehicle.ID.Hex())
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "rider detached from motor vehicle successfully", "data": vehicle})
}

func DeleteMotorVehicle(c *gin.Context) {
	id := c.Param("id")
	log.Printf("DeleteMotorVehicle: request received for id=%s", id)
	if err := services.DeleteMotorVehicle(c.Request.Context(), id); err != nil {
		log.Printf("DeleteMotorVehicle: service error for id=%s: %v", id, err)
		if errors.Is(err, services.ErrMotorVehicleNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "motor vehicle not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	log.Printf("DeleteMotorVehicle: vehicle deleted: %s", id)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "motor vehicle deleted successfully"})
}
