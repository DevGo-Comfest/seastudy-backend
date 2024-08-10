package service

import (
	"fmt"
	"sea-study/api/models"
	"sea-study/constants"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Check if a user is already enrolled in a course
func IsUserEnrolled(db *gorm.DB, userID uuid.UUID, courseID int) (bool, error) {
	var enrollment models.Enrollment
	if err := db.Where("user_id = ? AND course_id = ?", userID, courseID).First(&enrollment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Enroll a user in a course
func EnrollUser(db *gorm.DB, userID string, courseID int) (*models.Enrollment, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidUserID)
	}

	enrolled, err := IsUserEnrolled(db, userUUID, courseID)
	if err != nil {
		return nil, err
	}
	if enrolled {
		return nil, fmt.Errorf(constants.ErrUserAlreadyEnrolled)
	}

	var course models.Course
	if err := db.Where("course_id = ?", courseID).First(&course).Error; err != nil {
		return nil, err
	}

	balance, err := GetUserBalance(db, userUUID)
	if err != nil {
		return nil, err
	}
	if balance < float64(course.Price) {
		return nil, fmt.Errorf(constants.ErrInsufficientBalance)
	}

	// Start a transaction
	tx := db.Begin()

	if err := UpdateUserBalance(tx, userUUID, -float64(course.Price)); err != nil {
		tx.Rollback()
		return nil, err
	}

	enrollment := &models.Enrollment{
		UserID:       userUUID,
		CourseID:     courseID,
		DateEnrolled: time.Now(),
	}
	if err := tx.Create(enrollment).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf(constants.ErrFailedToCreateEnrollment)
	}

	tx.Commit()

	return enrollment, nil
}
