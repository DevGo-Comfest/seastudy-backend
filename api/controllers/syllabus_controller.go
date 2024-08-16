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

type CreateSyllabusInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	CourseID    int    `json:"course_id" binding:"required"`
}

type UpdateSyllabusInput struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
}

func CreateSyllabus(c *gin.Context, db *gorm.DB) {
	var input CreateSyllabusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidInput})
		return
	}

	instructorID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
		return
	}

	instructorUUID, err := uuid.Parse(instructorID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidUserID})
		return
	}

	syllabus := models.Syllabus{
		Title:        input.Title,
		Description:  input.Description,
		InstructorID: instructorUUID,
		CourseID:     input.CourseID,
	}

	if err := service.CreateSyllabus(db, &syllabus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToCreateSyllabus})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Syllabus created successfully", "syllabus": syllabus})
}

func UpdateSyllabus(c *gin.Context, db *gorm.DB) {
	var input UpdateSyllabusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidInput})
		return
	}

	syllabusID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidSyllabusID})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
		return
	}

	instructorUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidUserID})
		return
	}

	updatedSyllabus := models.Syllabus{
		Title:        input.Title,
		Description:  input.Description,
		InstructorID: instructorUUID,
	}

	if err := service.UpdateSyllabus(db, syllabusID, &updatedSyllabus, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToUpdateSyllabus})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Syllabus updated successfully", "syllabus": updatedSyllabus})
}

func DeleteSyllabus(c *gin.Context, db *gorm.DB) {
	syllabusID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidSyllabusID})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
		return
	}

	if err := service.DeleteSyllabus(db, syllabusID, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToDeleteSyllabus})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Syllabus deleted successfully"})
}

func GetSyllabus(c *gin.Context, db *gorm.DB) {
	syllabusID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidSyllabusID})
		return
	}

	syllabus, err := service.GetSyllabus(db, syllabusID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToRetrieveSyllabus})
		return
	}

	c.JSON(http.StatusOK, gin.H{"syllabus": syllabus})
}
