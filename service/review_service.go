package service

import (
	"fmt"
	"sea-study/api/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseReviewResponse struct {
	CourseReviewID int       `json:"course_review_id"`
	CourseID       int       `json:"course_id"`
	FeedbackText   string    `json:"feedback_text"`
	UserName       string    `json:"user_name"`
	UserRole       string    `json:"user_role"`
	Rate           int       `json:"rate"`
}

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

	// Check if the user has already submitted a review for this course
	var existingReview models.CourseReview
	err = db.Where("user_id = ? AND course_id = ?", userUUID, courseID).First(&existingReview).Error
	if err == nil {
		return nil, fmt.Errorf("user has already submitted a review for this course")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
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


func GetCourseReviews(db *gorm.DB, courseID int) ([]CourseReviewResponse, error) {
	var reviewResponses []CourseReviewResponse

	err := db.Table("course_reviews").
		Select("course_reviews.course_review_id, course_reviews.course_id, course_reviews.feedback_text, course_reviews.user_id, course_reviews.created_at, course_reviews.rate, users.name as user_name, users.role as user_role").
		Joins("left join users on course_reviews.user_id = users.user_id").
		Where("course_reviews.course_id = ?", courseID).
		Order("course_reviews.created_at desc").
		Scan(&reviewResponses).Error

	if err != nil {
		return nil, err
	}

	return reviewResponses, nil
}