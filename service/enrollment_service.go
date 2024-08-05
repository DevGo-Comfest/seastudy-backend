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

	enrollment := &models.Enrollment{
		UserID:       userUUID,
		CourseID:     courseID,
		DateEnrolled: time.Now(),
	}
	if err := db.Create(enrollment).Error; err != nil {
		return nil, err
	}
	return enrollment, nil
}
