package routes

import (
	"practice/api/shared/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api", middleware.SaveActivitiesInDB())
	RegisterUser(api)          // User Routes
	RegisterMobile(api)        // Mobile Routes
	RegisterShop(api)           // Shop Routes
	RegisterLocation(api)       // Location Routes
	RegisterMotorVehicles(api)  // Motor vehicle and registration routes
	RegisterCustomerRoutes(api) // Customer and ride routes
	RegisterState(api)          // State Routes
	RegisterCountry(api)        // Country Routes
	RegisterStage(api)          // Stage Routes
	RegisterGallery(api)        // Gallery Routes
	RegisterFruit(api)          // Fruit Routes
}
