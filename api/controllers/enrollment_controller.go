package controllers

import (
	"net/http"
	"sea-study/constants"
	"sea-study/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EnrollInput struct {
	CourseID int `json:"course_id" binding:"required"`
}

func EnrollUser(c *gin.Context, db *gorm.DB) {
	var input EnrollInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidInput})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
		return
	}

	enrollment, err := service.EnrollUser(db, userID.(string), input.CourseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "enrollment successful", "enrollment": enrollment})
}


func GetEnrolledCourses(c *gin.Context, db *gorm.DB) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
		return
	}

	courses, err := service.GetEnrolledCourses(db, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"enrolled_courses": courses})
}