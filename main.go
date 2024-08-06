package main

import (
	"sea-study/api/routes"
	"sea-study/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Initialize the database
	db := config.InitDB()

    // Enable CORS for all origins
    r.Use(cors.Default())
	
	if db != nil {
		r.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "DB connected successfully",
			})
		})
	} else {
		r.GET("/", func(c *gin.Context) {
			c.JSON(500, gin.H{
				"message": "Failed to connect to DB",
			})
		})
	}
	routes.RegisterUserRoutes(r,db)
	routes.RegisterEnrollmentRoutes(r, db)
	routes.RegisterTopupRoutes(r, db)
	routes.RegisterReviewRoutes(r, db)
	routes.RegisterForumPostRoutes(r, db)

	r.Run() // listen and serve on 0.0.0.0:8080
}
