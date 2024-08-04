package routes

import (
	"sea-study/api/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(router *gin.Engine, db *gorm.DB) {
    api := router.Group("/api")
    {
        api.POST("/register", func(c *gin.Context) {
            controllers.RegisterUser(c, db)
        })
    }
}
