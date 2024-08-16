package service

import (
	"fmt"
	"sea-study/api/models"
	"sea-study/constants"
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
        return nil, fmt.Errorf(constants.ErrInvalidUserID)
    }

    enrolled, err := IsUserEnrolled(db, userUUID, courseID)
    if err != nil {
        return nil, err
    }
    if !enrolled {
        return nil, fmt.Errorf(constants.ErrUserNotEnrolledInCourse)
    }

    // Check if the user has completed all syllabus materials
    progressPercentage, err := GetUserCourseProgress(db, userID, courseID)
    if err != nil {
        return nil, err
    }
    if progressPercentage < 100 {
        return nil, fmt.Errorf(constants.ErrIncompleteCourseProgress)
    }

    if rate < 1 || rate > 5 {
        return nil, fmt.Errorf(constants.ErrInvalidRate)
    }

    var existingReview models.CourseReview
    err = db.Where("user_id = ? AND course_id = ?", userUUID, courseID).First(&existingReview).Error
    if err == nil {
        return nil, fmt.Errorf(constants.ErrUserAlreadySubmittedReview)
    } else if err != gorm.ErrRecordNotFound {
        return nil, err
    }

    // Start a transaction
    tx := db.Begin()

    review := &models.CourseReview{
        CourseID:     courseID,
        FeedbackText: feedback,
        UserID:       userUUID,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
        Rate:         rate,
    }

    if err := tx.Create(review).Error; err != nil {
        tx.Rollback()
        return nil, fmt.Errorf(constants.ErrFailedToCreateReview)
    }

    // Update the course rating
    if err := updateCourseRating(tx, courseID); err != nil {
        tx.Rollback()
        return nil, err
    }

    tx.Commit()

    return review, nil
}

func updateCourseRating(tx *gorm.DB, courseID int) error {
    var result struct {
        TotalRating int64
        ReviewCount int64
    }
    if err := tx.Model(&models.CourseReview{}).
        Where("course_id = ?", courseID).
        Select("COALESCE(SUM(rate), 0) as total_rating, COUNT(*) as review_count").
        Scan(&result).Error; err != nil {
        return err
    }

    if result.ReviewCount == 0 {
        return fmt.Errorf(constants.ErrNoReviewsFound)
    }

    averageRating := int(result.TotalRating / result.ReviewCount)

    if err := tx.Model(&models.Course{}).Where("course_id = ?", courseID).Update("rating", averageRating).Error; err != nil {
        return fmt.Errorf(constants.ErrFailedToUpdateRating)
    }

    return nil
}



func GetCourseReviews(db *gorm.DB, courseID int) ([]CourseReviewResponse, error) {
	var reviewResponses []CourseReviewResponse

    err := db.Table("course_reviews").
		Select("course_reviews.course_review_id, course_reviews.course_id, course_reviews.feedback_text, course_reviews.user_id, course_reviews.created_at, course_reviews.rate, users.name as user_name, users.role as user_role").
		Joins("left join users on course_reviews.user_id = users.user_id").
		Joins("left join courses on course_reviews.course_id = courses.course_id").
		Where("course_reviews.course_id = ? AND courses.status = ?", courseID, models.ActiveStatus).
		Order("course_reviews.created_at desc").
		Scan(&reviewResponses).Error
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return reviewResponses, nil
}
