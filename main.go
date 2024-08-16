package main

import (
	"sea-study/api/routes"
	"sea-study/config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Static("/uploads", "./uploads")

	// Initialize the database
	db := config.InitDB()


	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},  
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: false,  
		MaxAge:           24 * time.Hour, 
	}))
	
	
	
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
	routes.RegisterCourseRoutes(r, db)
	routes.RegisterSyllabusRoutes(r, db)
	routes.RegisterAssignmentRoutes(r, db)
	routes.RegisterSyllabusMaterialRoutes(r, db)	
	routes.RegisterSubmissionRoutes(r, db)
	routes.RegisterProgressRoutes(r, db)


	r.Run() // listen and serve on 0.0.0.0:8080
}
