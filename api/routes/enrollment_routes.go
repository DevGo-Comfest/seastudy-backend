package routes

import (
	"sea-study/api/controllers"
	"sea-study/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterEnrollmentRoutes(router *gin.Engine, db *gorm.DB) {
    api := router.Group("/api")
    
    authenticated := api.Group("/")
    authenticated.Use(middleware.UserMiddleware())
    {
        authenticated.POST("/enroll", func(c *gin.Context) {
            controllers.EnrollUser(c, db)
        })
    }
}
