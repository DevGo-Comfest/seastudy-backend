package controllers

import (
	"net/http"
	"sea-study/api/models"
	"sea-study/constants"
	"sea-study/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OpenAssignmentInput struct {
	AssignmentID int `json:"assignment_id" binding:"required"`
}

func OpenAssignment(c *gin.Context, db *gorm.DB) {
	var input OpenAssignmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidInput})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
		return
	}

	err := service.OpenAssignment(db, userID.(string), input.AssignmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User assignment opened successfully"})
}

type AssignmentInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	MaximumTime int    `json:"maximum_time" binding:"required"`
}

func CreateAssignment(c *gin.Context, db *gorm.DB) {
	// Get the syllabus ID from the URL param
	syllabusIDParam := c.Param("id")
	syllabusID, err := strconv.Atoi(syllabusIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidSyllabusID})
		return
	}

	// Get the user ID from the middleware
	instructorID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
		return
	}
	userUUID, err := uuid.Parse(instructorID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidUserID})
		return
	}

	// Check if user is the instructor of the syllabus
	syllabus, err := service.GetSyllabus(db, syllabusID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if syllabus.InstructorID.String() != userUUID.String() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorizedAssignmentAction})
		return
	}

	var input AssignmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidInput})
		return
	}

	assignment := models.Assignment{
		SyllabusID:  syllabusID,
		Title:       input.Title,
		Description: input.Description,
		MaximumTime: input.MaximumTime,
	}

	createdAssignment, err := service.CreateAssignment(db, &assignment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToCreateAssignment})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assignment created successfully", "assignment": createdAssignment})
}

func GetAssignment(c *gin.Context, db *gorm.DB) {
	assignmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidAssignmentID})
		return
	}

	assignment, err := service.GetAssignmentByID(db, assignmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrAssignmentNotFound})
		return
	}

	// Check if the user is the instructor of the syllabus
	userID, _ := c.Get("userID")
	syllabus, err := service.GetSyllabus(db, assignment.SyllabusID)
	if err != nil || syllabus.InstructorID.String() != userID.(string) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorizedAssignmentAction})
		return
	}

	c.JSON(http.StatusOK, assignment)
}

func UpdateAssignment(c *gin.Context, db *gorm.DB) {
	assignmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidAssignmentID})
		return
	}

	var input AssignmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidInput})
		return
	}

	assignment, err := service.GetAssignmentByID(db, assignmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrAssignmentNotFound})
		return
	}

	// Check if the user is the instructor of the syllabus
	userID, _ := c.Get("userID")
	syllabus, err := service.GetSyllabus(db, assignment.SyllabusID)
	if err != nil || syllabus.InstructorID.String() != userID.(string) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorizedAssignmentAction})
		return
	}

	assignment.Title = input.Title
	assignment.Description = input.Description
	assignment.MaximumTime = input.MaximumTime

	updatedAssignment, err := service.UpdateAssignment(db, assignment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToUpdateAssignment})
		return
	}

	c.JSON(http.StatusOK, updatedAssignment)
}

func DeleteAssignment(c *gin.Context, db *gorm.DB) {
	assignmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidAssignmentID})
		return
	}

	assignment, err := service.GetAssignmentByID(db, assignmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrAssignmentNotFound})
		return
	}

	// Check if the user is the instructor of the syllabus
	userID, _ := c.Get("userID")
	syllabus, err := service.GetSyllabus(db, assignment.SyllabusID)
	if err != nil || syllabus.InstructorID.String() != userID.(string) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorizedAssignmentAction})
		return
	}

	if err := service.DeleteAssignment(db, assignmentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToDeleteAssignment})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assignment deleted successfully"})
}

func GetUserAssignment(c *gin.Context, db *gorm.DB) {
    assignmentID, err := strconv.Atoi(c.Param("assignmentId"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidAssignmentID})
        return
    }

    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
        return
    }

    userAssignment, err := service.GetUserAssignment(db, assignmentID, userID.(string))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToRetrieveUserAssignment})
        return
    }

    if userAssignment == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrUserAssignmentNotFound})
        return
    }

    c.JSON(http.StatusOK, userAssignment)
}