package routes

import (
	"practice/api/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterFlowerShop registers all flower shop CRUD routes plus the endpoints
// used to add and remove flowers from a shop's flowers array.
func RegisterFlowerShop(rg *gin.RouterGroup) {
	shops := rg.Group("/flowershops")
	{
		shops.POST("", handlers.CreateFlowerShop)
		shops.GET("", handlers.GetFlowerShops)
		shops.GET("/:id", handlers.GetFlowerShop)
		shops.PUT("/:id", handlers.UpdateFlowerShop)
		shops.DELETE("/:id", handlers.DeleteFlowerShop)
		shops.POST("/:id/flowers", handlers.AddFlowerToShop)
		shops.DELETE("/:id/flowers/:flowerid", handlers.RemoveFlowerFromShop)
	}
}
