package routes

import (
	"sea-study/api/controllers"
	"sea-study/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterSyllabusRoutes(router *gin.Engine, db *gorm.DB) {
	api := router.Group("/api")
	{
		// Public routes
		api.GET("/syllabus/:id", func(c *gin.Context) {
			controllers.GetSyllabus(c, db)
		})

		// Author only
		authorRoutes := api.Group("/")
		authorRoutes.Use(middleware.UserMiddleware(), middleware.AuthorMiddleware())
		{
			authorRoutes.POST("/syllabus", func(c *gin.Context) {
				controllers.CreateSyllabus(c, db)
			})
			authorRoutes.PUT("/syllabus/:id", func(c *gin.Context) {
				controllers.UpdateSyllabus(c, db)
			})
			authorRoutes.DELETE("/syllabus/:id", func(c *gin.Context) {
				controllers.DeleteSyllabus(c, db)
			})
		}
	}

}
