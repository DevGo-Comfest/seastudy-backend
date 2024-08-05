package service

import (
	"fmt"
	"sea-study/api/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func EnrollUser(db *gorm.DB, userID string, courseID int) (*models.Enrollment, error) {
	// Parse the userID string into a uuid.UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	var existingEnrollment models.Enrollment
	if err := db.Where("user_id = ? AND course_id = ?", userUUID, courseID).First(&existingEnrollment).Error; err == nil {
		return nil, fmt.Errorf("user is already enrolled in the course")
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
		return nil, fmt.Errorf("insufficient balance to enroll in the course")
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
		return nil, err
	}

	tx.Commit()

	return enrollment, nil
}
