package controllers

import (
	"net/http"
	"sea-study/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ForumPostInput struct {
	CourseID int    `json:"course_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

func CreateForumPost(c *gin.Context, db *gorm.DB) {
	var input ForumPostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	forumPost, err := service.CreateForumPost(db, userID.(string), input.CourseID, input.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Forum post created successfully", "forum_post": forumPost})
}

func GetForumPosts(c *gin.Context, db *gorm.DB) {
	courseIDParam := c.Param("course_id")
	courseID, err := strconv.Atoi(courseIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	forumPosts, err := service.GetForumPosts(db, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"forum_posts": forumPosts})
}