package routes

import (
	"sea-study/api/controllers"
	"sea-study/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterReviewRoutes(router *gin.Engine, db *gorm.DB) {
	api := router.Group("/api")
	
	authenticated := api.Group("/")
	authenticated.Use(middleware.UserMiddleware())
	{
		authenticated.POST("/review", func(c *gin.Context) {
			controllers.CreateReview(c, db)
		})
		authenticated.GET("/review/:course_id", func(c *gin.Context) {
			controllers.GetCourseReviews(c, db)
		})
	}
}
