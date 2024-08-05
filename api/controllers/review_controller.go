package controllers

import (
	"net/http"
	"sea-study/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewInput struct {
	CourseID     int    `json:"course_id" binding:"required"`
	FeedbackText string `json:"feedback_text" binding:"required"`
	Rate         int    `json:"rate" binding:"required"`
}

func CreateReview(c *gin.Context, db *gorm.DB) {
	var input ReviewInput
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
