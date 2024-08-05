package service

import (
	"fmt"
	"sea-study/api/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateReview(db *gorm.DB, userID string, courseID int, feedback string, rate int) (*models.CourseReview, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	enrolled, err := IsUserEnrolled(db, userUUID, courseID)
	if err != nil {
		return nil, err
	}
	if !enrolled {
		return nil, fmt.Errorf("user is not enrolled in the course")
	}
	if rate < 1 || rate > 5 {
		return nil, fmt.Errorf("rate must be between 1 and 5")
	}

	review := &models.CourseReview{
		CourseID:     courseID,
		FeedbackText: feedback,
		UserID:       userUUID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Rate:         rate,
	}

	if err := db.Create(review).Error; err != nil {
		return nil, err
	}

	return review, nil
}
