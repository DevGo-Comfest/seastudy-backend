package service

import (
	"errors"
	"sea-study/api/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateSubmission(db *gorm.DB, submission *models.Submission) (*models.Submission, error) {
	err := db.Transaction(func(tx *gorm.DB) error {
		// Get the assignment based on assignment id in submission
		var assignment models.Assignment
		if err := tx.Where("assignment_id = ?", submission.AssignmentID).First(&assignment).Error; err != nil {
			return err
		}

		// Get the syllabus based on the assignment's syllabus_id
		var syllabus models.Syllabus
		if err := tx.Where("syllabus_id = ?", assignment.SyllabusID).First(&syllabus).Error; err != nil {
			return err
		}

		// Create the submission
		if err := tx.Create(submission).Error; err != nil {
			return err
		}

		// Update or create the user progress
		var progress models.UserProgress
		result := tx.Where("user_id = ? AND course_id = ? AND syllabus_id = ?",
			submission.UserID, syllabus.CourseID, syllabus.SyllabusID).First(&progress)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				// Create new progress if not found
				progress = models.UserProgress{
					UserID:       submission.UserID,
					CourseID:     syllabus.CourseID,
					SyllabusID:   syllabus.SyllabusID,
					Status:       models.ProgressStatusEnum(models.Completed),
					LastAccessed: time.Now(),
				}
				if err := tx.Create(&progress).Error; err != nil {
					return err
				}
			} else {
				return result.Error
			}
		} else {
			// Update existing progress
			if err := tx.Model(&progress).Updates(models.UserProgress{
				Status:       models.ProgressStatusEnum(models.Completed),
				LastAccessed: time.Now(),
			}).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
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