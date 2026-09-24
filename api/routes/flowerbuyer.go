package routes

import (
	"practice/api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterFlowerBuyer(rg *gin.RouterGroup) {
	buyers := rg.Group("/flower-buyers")
	buyers.POST("", handlers.CreateFlowerBuyer)
	buyers.GET("", handlers.GetFlowerBuyers)
	buyers.GET("/:id", handlers.GetFlowerBuyer)
	buyers.PUT("/:id", handlers.UpdateFlowerBuyer)
	buyers.DELETE("/:id", handlers.DeleteFlowerBuyer)
}
