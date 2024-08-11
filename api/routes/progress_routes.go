package routes

import (
	"sea-study/api/controllers"
	"sea-study/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterProgressRoutes(router *gin.Engine, db *gorm.DB) {
	api := router.Group("/api")
	authenticated := api.Group("/")

	authenticated.Use(middleware.UserMiddleware(db))
	{
		authenticated.POST("/progress/update", func(c *gin.Context) {
			controllers.UpdateUserProgress(c, db)
		})
		authenticated.GET("/progress/course/:course_id", func(c *gin.Context) {
			controllers.GetUserCourseProgress(c, db)
		})

		// Add this route for instructors to get student progress
		authenticated.GET("/progress/course/:course_id/students", func(c *gin.Context) {
			controllers.GetStudentsProgressForCourse(c, db)
		})
	}
	
}
