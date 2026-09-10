package routes

import (
	"practice/api/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterShop registers all shop related routes.
func RegisterShop(rg *gin.RouterGroup) {
	shops := rg.Group("/shops")
	{
		shops.POST("", handlers.CreateShop)
		shops.GET("", handlers.GetShops)
		shops.GET("/:id", handlers.GetShop)
		shops.PUT("/:id", handlers.UpdateShop)
		shops.DELETE("/:id", handlers.DeleteShop)
	}
}
