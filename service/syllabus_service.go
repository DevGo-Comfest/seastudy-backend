package service

import (
	"errors"
	"fmt"
	"sea-study/api/models"

	"gorm.io/gorm"
)

// CreateSyllabus creates a new syllabus with automatic order assignment
func CreateSyllabus(db *gorm.DB, syllabus *models.Syllabus) error {
	// Find the highest order for the course
	var maxOrder int
	err := db.Model(&models.Syllabus{}).Where("course_id = ?", syllabus.CourseID).Select("COALESCE(MAX(\"order\"), 0)").Row().Scan(&maxOrder)
	if err != nil {
		return err
	}
	fmt.Printf("Max order for course_id %d: %d\n", syllabus.CourseID, maxOrder)
	
	// Assign the next order value
	syllabus.Order = maxOrder + 1
	fmt.Printf("Assigned order: %d\n", syllabus.Order)
	
	return db.Create(syllabus).Error
}

// UpdateSyllabus updates an existing syllabus with ownership check
func UpdateSyllabus(db *gorm.DB, syllabusID int, updatedSyllabus *models.Syllabus, userID string) error {
	var syllabus models.Syllabus
	if err := db.First(&syllabus, syllabusID).Error; err != nil {
		return err
	}

	if syllabus.InstructorID.String() != userID {
		return errors.New("unauthorized to update this syllabus")
	}

	// Only update allowed fields
	return db.Model(&syllabus).Updates(map[string]interface{}{
		"title":        updatedSyllabus.Title,
		"description":  updatedSyllabus.Description,
		"assignment_id": updatedSyllabus.AssignmentID,
	}).Error
}


// DeleteSyllabus deletes a syllabus by ID with ownership check and reorders the remaining items
func DeleteSyllabus(db *gorm.DB, syllabusID int, userID string) error {
	var syllabus models.Syllabus
	if err := db.First(&syllabus, syllabusID).Error; err != nil {
		return err
	}

	if syllabus.InstructorID.String() != userID {
		return errors.New("unauthorized to delete this syllabus")
	}

	// Begin a transaction
	tx := db.Begin()

	// Delete the syllabus
	if err := tx.Delete(&models.Syllabus{}, syllabusID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Reorder the remaining syllabuses
	if err := tx.Model(&models.Syllabus{}).Where("course_id = ? AND \"order\" > ?", syllabus.CourseID, syllabus.Order).Update("\"order\"", gorm.Expr("\"order\" - 1")).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	tx.Commit()

	return nil
}


// GetSyllabus gets a syllabus by ID
func GetSyllabus(db *gorm.DB, syllabusID int) (*models.Syllabus, error) {
	var syllabus models.Syllabus
	if err := db.Preload("Materials").Preload("Assignments").First(&syllabus, syllabusID).Error; err != nil {
		return nil, err
	}
	return &syllabus, nil
}
