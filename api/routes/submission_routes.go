package routes

import (
	"sea-study/api/controllers"
	"sea-study/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterSubmissionRoutes(router *gin.Engine, db *gorm.DB) {
	api := router.Group("/api")
	{
		authenticated := api.Group("/")
		authenticated.Use(middleware.UserMiddleware())
		{
			authenticated.POST("/assignments/:assignment_id/submissions", func(c *gin.Context) {
				controllers.CreateSubmission(c, db)
			})
			authenticated.PUT("/submissions/:id", func(c *gin.Context) {
				controllers.UpdateSubmission(c, db)
			})
			authenticated.DELETE("/submissions/:id", func(c *gin.Context) {
				controllers.DeleteSubmission(c, db)
			})
		}

		authorRoutes := api.Group("/")
		authorRoutes.Use(middleware.UserMiddleware(), middleware.AuthorMiddleware())
		{
			authorRoutes.PUT("/submissions/:id/grade", func(c *gin.Context) {
				controllers.GradeSubmission(c, db)
			})
		}
	}
}
