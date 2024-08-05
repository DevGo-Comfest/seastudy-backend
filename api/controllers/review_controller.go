package controllers

import (
	"net/http"
	"sea-study/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateReviewInput struct {
	CourseID     int    `json:"course_id" binding:"required"`
	FeedbackText string `json:"feedback_text" binding:"required"`
	Rate         int    `json:"rate" binding:"required"`
}

func CreateReview(c *gin.Context, db *gorm.DB) {
	var input CreateReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	review, err := service.CreateReview(db, userID.(string), input.CourseID, input.FeedbackText, input.Rate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review submitted successfully", "review": review})
}


func GetCourseReviews(c *gin.Context, db *gorm.DB) {
	courseIDParam := c.Param("course_id")
	courseID, err := strconv.Atoi(courseIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	reviews, err := service.GetCourseReviews(db, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}