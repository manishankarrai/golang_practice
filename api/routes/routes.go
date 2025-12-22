package routes

import (
	"practice/api/shared/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api", middleware.SaveActivitiesInDB())
	RegisterUser(api) // User Routes
}
