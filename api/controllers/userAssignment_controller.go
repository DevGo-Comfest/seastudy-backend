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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid syllabus ID"})
		return
	}

	// Get the user ID from the middleware
	instructorID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userUUID, err := uuid.Parse(instructorID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user is the instructor of the syllabus
	syllabus, err := service.GetSyllabus(db, syllabusID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if syllabus.InstructorID.String() != userUUID.String() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input AssignmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assignment created successfully", "assignment": createdAssignment})
}

func GetAssignment(c *gin.Context, db *gorm.DB) {
	assignmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	assignment, err := service.GetAssignmentByID(db, assignmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	// Check if the user is the instructor of the syllabus
	userID, _ := c.Get("userID")
	syllabus, err := service.GetSyllabus(db, assignment.SyllabusID)
	if err != nil || syllabus.InstructorID.String() != userID.(string) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, assignment)
}

func UpdateAssignment(c *gin.Context, db *gorm.DB) {
	assignmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	var input AssignmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assignment, err := service.GetAssignmentByID(db, assignmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	// Check if the user is the instructor of the syllabus
	userID, _ := c.Get("userID")
	syllabus, err := service.GetSyllabus(db, assignment.SyllabusID)
	if err != nil || syllabus.InstructorID.String() != userID.(string) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	assignment.Title = input.Title
	assignment.Description = input.Description
	assignment.MaximumTime = input.MaximumTime

	updatedAssignment, err := service.UpdateAssignment(db, assignment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedAssignment)
}

func DeleteAssignment(c *gin.Context, db *gorm.DB) {
	assignmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	assignment, err := service.GetAssignmentByID(db, assignmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	// Check if the user is the instructor of the syllabus
	userID, _ := c.Get("userID")
	syllabus, err := service.GetSyllabus(db, assignment.SyllabusID)
	if err != nil || syllabus.InstructorID.String() != userID.(string) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := service.DeleteAssignment(db, assignmentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assignment deleted successfully"})
}