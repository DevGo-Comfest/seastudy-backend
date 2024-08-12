package service

import (
	"errors"
	"fmt"
	"sea-study/api/models"
	"sea-study/constants"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateSyllabus creates a new syllabus with automatic order assignment
func CreateSyllabus(db *gorm.DB, syllabus *models.Syllabus, userID uuid.UUID) error {
	// Check if the user is the primary author or an instructor for the course
	var count int64
	err := db.Table("course_instructors").
		Where("course_instructors.course_id = ? AND course_instructors.instructor_id = ?", syllabus.CourseID, userID).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count == 0 {
		err = db.Model(&models.Course{}).
			Where("course_id = ? AND primary_author = ?", syllabus.CourseID, userID).
			Count(&count).Error

		if err != nil || count == 0 {
			return fmt.Errorf(constants.ErrUnauthorized)
		}
	}

	var maxOrder int
	err = db.Model(&models.Syllabus{}).Where("course_id = ?", syllabus.CourseID).Select("COALESCE(MAX(\"order\"), 0)").Row().Scan(&maxOrder)
	if err != nil {
		return err
	}

	syllabus.Order = maxOrder + 1

	return db.Create(syllabus).Error
}

// UpdateSyllabus updates an existing syllabus with ownership check
func UpdateSyllabus(db *gorm.DB, syllabusID int, updatedSyllabus *models.Syllabus, userID string) error {
	var syllabus models.Syllabus
	if err := db.First(&syllabus, syllabusID).Error; err != nil {
		return err
	}

	if syllabus.InstructorID.String() != userID {
		return errors.New(constants.ErrUnauthorizedSyllabus)
	}

	// Only update allowed fields
	return db.Model(&syllabus).Updates(map[string]interface{}{
		"title":        updatedSyllabus.Title,
		"description":  updatedSyllabus.Description,
	}).Error
}


// DeleteSyllabus deletes a syllabus by ID with ownership check and reorders the remaining items
func DeleteSyllabus(db *gorm.DB, syllabusID int, userID string) error {
	var syllabus models.Syllabus
	if err := db.First(&syllabus, syllabusID).Error; err != nil {
		return err
	}

	if syllabus.InstructorID.String() != userID {
		return errors.New(constants.ErrUnauthorizedSyllabus)
	}

	// Begin a transaction
	tx := db.Begin()

	if err := tx.Delete(&models.Syllabus{}, syllabusID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Reorder the remaining syllabuses
	if err := tx.Model(&models.Syllabus{}).Where("course_id = ? AND \"order\" > ?", syllabus.CourseID, syllabus.Order).Update("\"order\"", gorm.Expr("\"order\" - 1")).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}


// GetSyllabus gets a syllabus by ID
func GetSyllabus(db *gorm.DB, syllabusID int) (*models.Syllabus, error) {
	var syllabus models.Syllabus
	if err := db.Preload("Materials").Preload("Assignments").First(&syllabus, syllabusID).Error; err != nil {
		return nil, errors.New(constants.ErrSyllabusNotFound)
	}
	return &syllabus, nil
}
