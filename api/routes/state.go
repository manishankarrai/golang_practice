package routes

import (
	"practice/api/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterState registers all state related routes.
func RegisterState(rg *gin.RouterGroup) {
	states := rg.Group("/states")
	{
		states.POST("", handlers.CreateState)
		states.GET("", handlers.GetStates)
		states.GET("/:id", handlers.GetState)
		states.PUT("/:id", handlers.UpdateState)
		states.DELETE("/:id", handlers.DeleteState)
	}
}
