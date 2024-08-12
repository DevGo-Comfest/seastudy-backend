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
		api.GET("/courses/popular", func(c *gin.Context) {
			controllers.GetPopularCourses(c, db)
		})
		api.GET("/courses", func(c *gin.Context) {
			controllers.GetAllCourses(c, db)
		})
		api.GET("/courses/search", func(c *gin.Context) {
			controllers.SearchCourses(c, db)
		})

		// User only
		userRoutes := api.Group("/")
		userRoutes.Use(middleware.UserMiddleware(db))
		{
			userRoutes.GET("/courses/:course_id", func(c *gin.Context) {
				controllers.GetCourse(c, db)
			})
		}

		// Author only
		authorRoutes := api.Group("/")
		authorRoutes.Use(middleware.UserMiddleware(db), middleware.AuthorMiddleware())
		{
			authorRoutes.POST("/courses", func(c *gin.Context) {
				controllers.CreateCourse(c, db)
			})
			authorRoutes.GET("/courses/me", func(c *gin.Context) {
				controllers.GetMyCourse(c, db)
			})
			authorRoutes.POST("/courses/:course_id/instructors", func(c *gin.Context) {
				controllers.AddCourseInstructors(c, db)
			})
			authorRoutes.PUT("/courses/:course_id", func(c *gin.Context) {
				controllers.UpdateCourse(c, db)
			})
			authorRoutes.DELETE("/courses/:course_id", func(c *gin.Context) {
				controllers.DeleteCourse(c, db)
			})
			authorRoutes.POST("/courses/upload/image", func(c *gin.Context) {
				controllers.UploadCourseImage(c, db)
			})
			authorRoutes.PUT("/courses/:course_id/activate", func(c *gin.Context) {
				controllers.ActivateCourse(c, db)
			})
		}
	}
}
