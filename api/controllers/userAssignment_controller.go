package controllers

import (
	"net/http"
	"sea-study/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OpenAssignmentInput struct {
	AssignmentID int `json:"assignment_id" binding:"required"`
}

func OpenAssignment(c *gin.Context, db *gorm.DB) {
	var input OpenAssignmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := service.OpenAssignment(db, userID.(string), input.AssignmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User assignment opened successfully"})
}
