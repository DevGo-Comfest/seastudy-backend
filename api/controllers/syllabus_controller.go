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

type CreateSyllabusInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	AssignmentID string `json:"assignment_id"`
	CourseID    int    `json:"course_id" binding:"required"`
}

type UpdateSyllabusInput struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	AssignmentID string `json:"assignment_id"`
}

func CreateSyllabus(c *gin.Context, db *gorm.DB) {
	var input CreateSyllabusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	instructorID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	instructorUUID, err := uuid.Parse(instructorID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID format"})
		return
	}

	syllabus := models.Syllabus{
		Title:        input.Title,
		Description:  input.Description,
		InstructorID: instructorUUID,
		// AssignmentID: input.AssignmentID,
		CourseID:     input.CourseID,
	}

	if err := service.CreateSyllabus(db, &syllabus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Syllabus created successfully", "syllabus": syllabus})
}

func UpdateSyllabus(c *gin.Context, db *gorm.DB) {
	var input UpdateSyllabusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	syllabusID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid syllabus ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse the userID string to uuid.UUID
	instructorUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID format"})
		return
	}

	updatedSyllabus := models.Syllabus{
		Title:        input.Title,
		Description:  input.Description,
		InstructorID: instructorUUID, // Use the parsed uuid.UUID value here
		// AssignmentID: input.AssignmentID,
	}

	if err := service.UpdateSyllabus(db, syllabusID, &updatedSyllabus, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Syllabus updated successfully", "syllabus": updatedSyllabus})
}


func DeleteSyllabus(c *gin.Context, db *gorm.DB) {
	syllabusID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid syllabus ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := service.DeleteSyllabus(db, syllabusID, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Syllabus deleted successfully"})
}

func GetSyllabus(c *gin.Context, db *gorm.DB) {
	syllabusID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid syllabus ID"})
		return
	}

	syllabus, err := service.GetSyllabus(db, syllabusID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"syllabus": syllabus})
}
