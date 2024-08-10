package controllers

import (
	"net/http"
	"sea-study/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UpdateProgressInput struct {
	CourseID   int `json:"course_id" binding:"required"`
	SyllabusID int `json:"syllabus_id" binding:"required"`
}

func UpdateUserProgress(c *gin.Context, db *gorm.DB) {
	var input UpdateProgressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := service.UpdateUserProgress(db, userID.(string), input.CourseID, input.SyllabusID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User progress updated successfully"})
}

func GetUserCourseProgress(c *gin.Context, db *gorm.DB) {
	courseID, err := strconv.Atoi(c.Param("course_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	progressPercentage, err := service.GetUserCourseProgress(db, userID.(string), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"progress": progressPercentage})
}