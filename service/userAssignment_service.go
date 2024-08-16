package service

import (
	"errors"
	"fmt"
	"sea-study/api/models"
	"sea-study/constants"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateAssignment(db *gorm.DB, assignment *models.Assignment) (*models.Assignment, error) {
	if err := db.Create(assignment).Error; err != nil {
		return nil, fmt.Errorf(constants.ErrFailedToCreateAssignment)
	}
	return assignment, nil
}

func OpenAssignment(db *gorm.DB, userID string, assignmentID int) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf(constants.ErrInvalidUserID)
	}

	assignment, err := GetAssignmentByID(db, assignmentID)
	if err != nil {
		return fmt.Errorf(constants.ErrFailedToRetrieveAssignment)
	}

	dueDate := time.Now().AddDate(0, 0, assignment.MaximumTime)

	userAssignment := models.UserAssignment{
		AssignmentID: assignmentID,
		UserID:       userUUID,
		DueDate:      dueDate,
		CreatedAt:    time.Now(),
	}

	if err := CreateUserAssignment(db, &userAssignment); err != nil {
		return fmt.Errorf(constants.ErrFailedToCreateUserAssignment)
	}

	return nil
}

func GetAssignmentByID(db *gorm.DB, assignmentID int) (*models.Assignment, error) {
	var assignment models.Assignment
	result := db.Preload("Submissions").First(&assignment, assignmentID)
	if result.Error != nil {
		return nil, fmt.Errorf(constants.ErrFailedToRetrieveAssignment)
	}
	return &assignment, nil
}

func CreateUserAssignment(db *gorm.DB, userAssignment *models.UserAssignment) error {
	return db.Create(userAssignment).Error
}

func UpdateAssignment(db *gorm.DB, assignment *models.Assignment) (*models.Assignment, error) {
	if err := db.Save(assignment).Error; err != nil {
		return nil, fmt.Errorf(constants.ErrFailedToUpdateAssignment)
	}
	return assignment, nil
}

func DeleteAssignment(db *gorm.DB, assignmentID int) error {
	if err := db.Delete(&models.Assignment{}, assignmentID).Error; err != nil {
		return fmt.Errorf(constants.ErrFailedToDeleteAssignment)
	}
	return nil
}

func GetUserAssignment(db *gorm.DB, assignmentID int, userID string) (*models.UserAssignment, error) {
    var userAssignment models.UserAssignment
    
    result := db.Where("assignment_id = ? AND user_id = ?", assignmentID, userID).First(&userAssignment)
    
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, nil 
        }
        return nil, result.Error 
    }
    
    return &userAssignment, nil
}