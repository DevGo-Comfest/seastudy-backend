package routes

import (
	"sea-study/api/controllers"
	"sea-study/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAssignmentRoutes(router *gin.Engine, db *gorm.DB) {
	api := router.Group("/api")
	{
		authenticated := api.Group("/")
		authenticated.Use(middleware.UserMiddleware())
		{

			authenticated.POST("/assignments/open", func(c *gin.Context) {
				controllers.OpenAssignment(c, db)
			})
			authenticated.GET("/assignments/:id", func(c *gin.Context) {
				controllers.GetAssignment(c, db)
			})
		}

		authorRoutes := api.Group("/")
		authorRoutes.Use(middleware.UserMiddleware(), middleware.AuthorMiddleware())
		{
			authorRoutes.POST("/syllabus/:id/assignments", func(c *gin.Context) {
				controllers.CreateAssignment(c, db)
			})
			authorRoutes.PUT("/assignments/:id", func(c *gin.Context) {
				controllers.UpdateAssignment(c, db)
			})
			authorRoutes.DELETE("/assignments/:id", func(c *gin.Context) {
				controllers.DeleteAssignment(c, db)
			})
		}

	}

}
