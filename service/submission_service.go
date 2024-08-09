package service

import (
	"errors"
	"sea-study/api/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateSubmission(db *gorm.DB, submission *models.Submission) (*models.Submission, error) {
	if err := db.Create(submission).Error; err != nil {
		return nil, err
	}
	return submission, nil
}

func UpdateSubmission(db *gorm.DB, submission *models.Submission) (*models.Submission, error) {
	if err := db.Save(submission).Error; err != nil {
		return nil, err
	}
	return submission, nil
}

func DeleteSubmission(db *gorm.DB, submissionID int) error {
	result := db.Delete(&models.Submission{}, submissionID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("submission not found")
	}
	return nil
}

func GetSubmissionByID(db *gorm.DB, submissionID int) (*models.Submission, error) {
	var submission models.Submission
	if err := db.First(&submission, submissionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("submission not found")
		}
		return nil, err
	}
	return &submission, nil
}

func GetSubmissionByUserAndAssignment(db *gorm.DB, userID uuid.UUID, assignmentID int) (*models.Submission, error) {
	var submission models.Submission
	result := db.Where("user_id = ? AND assignment_id = ?", userID, assignmentID).First(&submission)
	if result.Error != nil {
		// No submission found
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &submission, nil
}
