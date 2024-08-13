package service

import (
	"fmt"
	"sea-study/api/models"
	"sea-study/constants"

	"gorm.io/gorm"
)

// CreateSyllabusMaterial creates a new syllabus material with automatic order assignment
func CreateSyllabusMaterial(db *gorm.DB, material *models.SyllabusMaterial) error {
	// Find the highest order for the syllabus
	var maxOrder int
	err := db.Model(&models.SyllabusMaterial{}).Where("syllabus_id = ?", material.SyllabusID).Select("COALESCE(MAX(\"order\"), 0)").Row().Scan(&maxOrder)
	if err != nil {
		return fmt.Errorf(constants.ErrFailedToCreateSyllabusMaterial)
	}
	
	// Assign the next order value
	material.Order = maxOrder + 1
	return db.Create(material).Error
}

func UpdateSyllabusMaterial(db *gorm.DB, materialID int, updatedMaterial *models.SyllabusMaterial, userID string) error {
	var material models.SyllabusMaterial
	if err := db.First(&material, materialID).Error; err != nil {
		return fmt.Errorf(constants.ErrInvalidSyllabusMaterialID)
	}

	var syllabus models.Syllabus
	if err := db.First(&syllabus, material.SyllabusID).Error; err != nil {
		return fmt.Errorf(constants.ErrInvalidSyllabusID)
	}

	if syllabus.InstructorID.String() != userID {
		return fmt.Errorf(constants.ErrUnauthorizedSyllabusAction)
	}

	return db.Model(&material).Updates(map[string]interface{}{
		"title":        updatedMaterial.Title,
		"description":  updatedMaterial.Description,
		"url_material": updatedMaterial.URLMaterial,
		"time_needed":  updatedMaterial.TimeNeeded,
	}).Error
}

// DeleteSyllabusMaterial deletes a syllabus material by ID 
func DeleteSyllabusMaterial(db *gorm.DB, materialID int, userID string) error {
	var material models.SyllabusMaterial
	if err := db.First(&material, materialID).Error; err != nil {
		return fmt.Errorf(constants.ErrInvalidSyllabusMaterialID)
	}

	var syllabus models.Syllabus
	if err := db.First(&syllabus, material.SyllabusID).Error; err != nil {
		return fmt.Errorf(constants.ErrInvalidSyllabusID)
	}

	// Begin a transaction
	tx := db.Begin()

	if err := tx.Delete(&models.SyllabusMaterial{}, materialID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf(constants.ErrFailedToDeleteSyllabusMaterial)
	}

	// Reorder the remaining materials
	if err := tx.Model(&models.SyllabusMaterial{}).Where("syllabus_id = ? AND \"order\" > ?", material.SyllabusID, material.Order).Update("\"order\"", gorm.Expr("\"order\" - 1")).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf(constants.ErrFailedToUpdateSyllabusMaterial)
	}

	tx.Commit()

	return nil
}

// GetSyllabusMaterial gets a syllabus material by ID
func GetSyllabusMaterial(db *gorm.DB, materialID int) (*models.SyllabusMaterial, error) {
	var material models.SyllabusMaterial
	if err := db.First(&material, materialID).Error; err != nil {
		return nil, fmt.Errorf(constants.ErrFailedToRetrieveSyllabusMaterial)
	}
	return &material, nil
}

// GetSyllabusMaterialsBySyllabusID gets all syllabus materials for a syllabus
func GetSyllabusMaterialsBySyllabusID(db *gorm.DB, syllabusID int) ([]models.SyllabusMaterial, error) {
	var materials []models.SyllabusMaterial
	if err := db.Where("syllabus_id = ?", syllabusID).Order("\"order\"").Find(&materials).Error; err != nil {
		return nil, fmt.Errorf(constants.ErrFailedToRetrieveSyllabusMaterial)
	}
	return materials, nil
}