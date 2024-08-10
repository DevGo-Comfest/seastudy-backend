package routes

import (
	"sea-study/api/controllers"
	"sea-study/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterForumPostRoutes(router *gin.Engine, db *gorm.DB) {
	api := router.Group("/api")

	authenticated := api.Group("/")
	authenticated.Use(middleware.UserMiddleware(db))
	{
		authenticated.POST("/forum-post", func(c *gin.Context) {
			controllers.CreateForumPost(c, db)
		})
		authenticated.GET("/forum-post/:course_id", func(c *gin.Context) {
			controllers.GetForumPosts(c, db)
		})
	}
}
