package controllers

import (
	"net/http"
	"sea-study/api/models"
	"sea-study/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubmissionInput struct {
	Status     models.SubmissionStatusEnum `json:"status"`
	Grade      int                         `json:"grade"`
	ContentURL string                      `json:"content_url" binding:"required"`
	IsLate     bool                        `json:"is_late"`
}

func CreateSubmission(c *gin.Context, db *gorm.DB) {
	// Get the assignment ID from the URL param
	assignmentIDParam := c.Param("assignment_id")
	assignmentID, err := strconv.Atoi(assignmentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	// Get the user ID from the middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if a submission already exists for this user and assignment
	existingSubmission, err := service.GetSubmissionByUserAndAssignment(db, userUUID, assignmentID)
	if err == nil && existingSubmission != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "You have already submitted for this assignment"})
		return
	}

	// If no existing submission is found, proceed with creating a new one
	var input SubmissionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	submission := models.Submission{
		Status:       input.Status,
		Grade:        input.Grade,
		ContentURL:   input.ContentURL,
		IsLate:       input.IsLate,
		AssignmentID: assignmentID,
		UserID:       userUUID,
	}

	createdSubmission, err := service.CreateSubmission(db, &submission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Submission created successfully", "submission": createdSubmission})
}

func UpdateSubmission(c *gin.Context, db *gorm.DB) {
    submissionID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
        return
    }

    var input SubmissionInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    submission, err := service.GetSubmissionByID(db, submissionID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
        return
    }

    // Check if the user is the owner of the submission
    userID, _ := c.Get("userID")
    if submission.UserID.String() != userID.(string) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    submission.Status = input.Status
    submission.Grade = input.Grade
    submission.ContentURL = input.ContentURL
    submission.IsLate = input.IsLate

    updatedSubmission, err := service.UpdateSubmission(db, submission)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, updatedSubmission)
}

func DeleteSubmission(c *gin.Context, db *gorm.DB) {
    submissionID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
        return
    }

    submission, err := service.GetSubmissionByID(db, submissionID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
        return
    }

    // Check if the user is the owner of the submission
    userID, _ := c.Get("userID")
    if submission.UserID.String() != userID.(string) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    if err := service.DeleteSubmission(db, submissionID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Submission deleted successfully"})
}

func GradeSubmission(c *gin.Context, db *gorm.DB) {
    submissionID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
        return
    }

    var input struct {
        Grade int `json:"grade" binding:"required"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    submission, err := service.GetSubmissionByID(db, submissionID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
        return
    }

    // Check if the user is an instructor
    userRole, _ := c.Get("userRole")
    if userRole != "author" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    submission.Grade = input.Grade

    updatedSubmission, err := service.UpdateSubmission(db, submission)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, updatedSubmission)
}

