package routes

import (
	"practice/api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterCustomerRoutes(rg *gin.RouterGroup) {
	customers := rg.Group("/customers")
	{
		customers.POST("", handlers.CreateCustomer)
		customers.GET("", handlers.GetCustomers)
		customers.GET("/:id", handlers.GetCustomer)
		customers.PUT("/:id", handlers.UpdateCustomer)
		customers.DELETE("/:id", handlers.DeleteCustomer)
		customers.PUT("/:id/location", handlers.UpdateCustomerLocation)
		customers.POST("/:id/rides", handlers.BookRide)
	}

	rides := rg.Group("/rides")
	{
		rides.GET("/available-vehicles", handlers.GetAvailableVehicles)
		rides.GET("", handlers.GetRides)
		rides.GET("/:id", handlers.GetRide)
	}
}
