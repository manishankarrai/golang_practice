package routes

import (
	"practice/api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterShopOrder(rg *gin.RouterGroup) {
	orders := rg.Group("/shop-orders")
	orders.POST("", handlers.CreateShopOrder)
	orders.GET("", handlers.GetShopOrders)
	orders.GET("/:id", handlers.GetShopOrder)
	orders.PUT("/:id", handlers.UpdateShopOrder)
	orders.DELETE("/:id", handlers.DeleteShopOrder)
}
