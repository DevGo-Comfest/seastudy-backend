package controllers

import (
	"net/http"
	"sea-study/api/models"
	"sea-study/constants"
	"sea-study/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateSyllabusMaterialInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	URLMaterial string `json:"url_material" binding:"required"`
	TimeNeeded  string `json:"time_needed"`
	SyllabusID  int    `json:"syllabus_id" binding:"required"`
}

type UpdateSyllabusMaterialInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URLMaterial string `json:"url_material"`
	TimeNeeded  string `json:"time_needed"`
}

func CreateSyllabusMaterial(c *gin.Context, db *gorm.DB) {
	var input CreateSyllabusMaterialInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidInput})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
		return
	}

	// Check if the user is the instructor for the syllabus
	var syllabus models.Syllabus
	if err := db.First(&syllabus, input.SyllabusID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidSyllabusID})
		return
	}
	if syllabus.InstructorID.String() != userID.(string) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorizedSyllabusAction})
		return
	}

	material := models.SyllabusMaterial{
		Title:       input.Title,
		Description: input.Description,
		URLMaterial: input.URLMaterial,
		TimeNeeded:  input.TimeNeeded,
		SyllabusID:  input.SyllabusID,
	}

	if err := service.CreateSyllabusMaterial(db, &material); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToCreateSyllabusMaterial})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Syllabus material created successfully", "syllabus_material": material})
}

func UpdateSyllabusMaterial(c *gin.Context, db *gorm.DB) {
	var input UpdateSyllabusMaterialInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidInput})
		return
	}

	materialID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidSyllabusMaterialID})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
		return
	}

	// Check if the user is the instructor for the syllabus
	var material models.SyllabusMaterial
	if err := db.First(&material, materialID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidSyllabusMaterialID})
		return
	}

	var syllabus models.Syllabus
	if err := db.First(&syllabus, material.SyllabusID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidSyllabusID})
		return
	}
	if syllabus.InstructorID.String() != userID.(string) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorizedSyllabusAction})
		return
	}
	updatedMaterial := models.SyllabusMaterial{
		Title:       input.Title,
		Description: input.Description,
		URLMaterial: input.URLMaterial,
		TimeNeeded:  input.TimeNeeded,
	}

	if err := service.UpdateSyllabusMaterial(db, materialID, &updatedMaterial, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToUpdateSyllabusMaterial})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Syllabus material updated successfully", "syllabus_material": updatedMaterial})
}

func DeleteSyllabusMaterial(c *gin.Context, db *gorm.DB) {
	materialID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidSyllabusMaterialID})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
		return
	}

	// Check if the user is the instructor for the syllabus
	var material models.SyllabusMaterial
	if err := db.First(&material, materialID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidSyllabusMaterialID})
		return
	}

	var syllabus models.Syllabus
	if err := db.First(&syllabus, material.SyllabusID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidSyllabusID})
		return
	}
	if syllabus.InstructorID.String() != userID.(string) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorizedSyllabusAction})
		return
	}

	if err := service.DeleteSyllabusMaterial(db, materialID, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToDeleteSyllabusMaterial})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Syllabus material deleted successfully"})
}

func GetSyllabusMaterial(c *gin.Context, db *gorm.DB) {
	materialID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidSyllabusMaterialID})
		return
	}

	material, err := service.GetSyllabusMaterial(db, materialID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToRetrieveSyllabusMaterial})
		return
	}

	c.JSON(http.StatusOK, gin.H{"syllabus_material": material})
}

func GetSyllabusMaterialBySyllabus(c *gin.Context, db *gorm.DB) {
	syllabusID, err := strconv.Atoi(c.Param("syllabus_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidSyllabusID})
		return
	}

	// Check if the user is authenticated or not
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
		return
	}

	materials, err := service.GetSyllabusMaterialsBySyllabusID(db, syllabusID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToRetrieveSyllabusMaterial})
		return
	}

	c.JSON(http.StatusOK, gin.H{"syllabus_materials": materials})
}
