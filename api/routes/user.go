package routes

import "github.com/gin-gonic/gin"

func RegisterUser(rg *gin.RouterGroup) {
	user := rg.Group("/user")
	{
		user.GET("/ping", func(c *gin.Context) {
			c.String(200, "text", "working fine")
		})
	}
}
