package routes

import (
	"sea-study/api/controllers"
	"sea-study/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterSyllabusMaterialRoutes(router *gin.Engine, db *gorm.DB) {
	api := router.Group("/api")
	{
		// Public routes
		api.GET("/syllabus-material/:id", func(c *gin.Context) {
			controllers.GetSyllabusMaterial(c, db)
		})

		// Author only
		authorRoutes := api.Group("/")
		authorRoutes.Use(middleware.UserMiddleware(), middleware.AuthorMiddleware())
		{
			authorRoutes.POST("/syllabus-material", func(c *gin.Context) {
				controllers.CreateSyllabusMaterial(c, db)
			})
			authorRoutes.PUT("/syllabus-material/:id", func(c *gin.Context) {
				controllers.UpdateSyllabusMaterial(c, db)
			})
			authorRoutes.DELETE("/syllabus-material/:id", func(c *gin.Context) {
				controllers.DeleteSyllabusMaterial(c, db)
			})
		}
	}
}
