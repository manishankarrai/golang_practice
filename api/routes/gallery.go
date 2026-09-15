package routes

import (
	"practice/api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterGallery(rg *gin.RouterGroup) {
	gallery := rg.Group("/gallery")
	{
		gallery.POST("", handlers.CreateGallery)
		gallery.GET("", handlers.GetGalleries)
		gallery.GET("/:id", handlers.GetGallery)
		gallery.PUT("/:id", handlers.UpdateGallery)
		gallery.DELETE("/:id", handlers.DeleteGallery)
	}
}
