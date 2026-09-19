package routes

import (
	"practice/api/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterMotorVehicles registers vehicle, location, rider, and owner-registration endpoints.
func RegisterMotorVehicles(rg *gin.RouterGroup) {
	vehicles := rg.Group("/motor-vehicles")
	{
		vehicles.POST("", handlers.CreateMotorVehicle)
		vehicles.GET("", handlers.GetMotorVehicles)
		vehicles.GET("/:id", handlers.GetMotorVehicle)
		vehicles.PUT("/:id", handlers.UpdateMotorVehicle)
		vehicles.DELETE("/:id", handlers.DeleteMotorVehicle)
		vehicles.POST("/:id/rider", handlers.AttachRiderToMotorVehicle)
		vehicles.DELETE("/:id/rider", handlers.DetachRiderFromMotorVehicle)
		vehicles.POST("/:id/location", handlers.UpdateVehicleLocation)
		vehicles.GET("/:id/location", handlers.GetVehicleLocation)
	}

	riders := rg.Group("/riders")
	{
		riders.POST("", handlers.CreateRider)
		riders.GET("", handlers.GetRiders)
		riders.GET("/:id", handlers.GetRider)
		riders.PUT("/:id", handlers.UpdateRider)
		riders.DELETE("/:id", handlers.DeleteRider)
	}

	registrations := rg.Group("/registrations")
	{
		registrations.POST("", handlers.CreateRegistration)
		registrations.GET("", handlers.GetRegistrations)
		registrations.GET("/:id", handlers.GetRegistration)
		registrations.PUT("/:id", handlers.UpdateRegistration)
		registrations.DELETE("/:id", handlers.DeleteRegistration)
	}
}
