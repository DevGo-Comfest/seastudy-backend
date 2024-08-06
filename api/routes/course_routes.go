package routes

import (
	"sea-study/api/controllers"
	"sea-study/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterCourseRoutes(router *gin.Engine, db *gorm.DB) {
	api := router.Group("/api")
	{
		// Public routes
		api.GET("/courses", func(c *gin.Context) {
			controllers.GetAllCourses(c, db)
		})
		api.GET("/courses/:course_id", func(c *gin.Context) {
			controllers.GetCourse(c, db)
		})

		// Author only
		authorRoutes := api.Group("/")
		authorRoutes.Use(middleware.UserMiddleware(), middleware.AuthorMiddleware())
		{
			authorRoutes.POST("/courses", func(c *gin.Context) {
				controllers.CreateCourse(c, db)
			})
			authorRoutes.PUT("/courses/:course_id", func(c *gin.Context) {
				controllers.UpdateCourse(c, db)
			})
			authorRoutes.DELETE("/courses/:course_id", func(c *gin.Context) {
				controllers.DeleteCourse(c, db)
			})
		}
	}
}
