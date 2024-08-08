package service

import (
	"fmt"
	"sea-study/api/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func OpenAssignment(db *gorm.DB, userID string, assignmentID int) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	assignment, err := GetAssignmentByID(db, assignmentID)
	if err != nil {
		return fmt.Errorf("failed to get assignment: %v", err)
	}

	dueDate := time.Now().AddDate(0, 0, assignment.MaximumTime)

	userAssignment := models.UserAssignment{
		AssignmentID: assignmentID,
		UserID:       userUUID,
		DueDate:      dueDate,
		CreatedAt:    time.Now(),
	}

	if err := CreateUserAssignment(db, &userAssignment); err != nil {
		return fmt.Errorf("failed to create user assignment: %v", err)
	}

	return nil
}

func GetAssignmentByID(db *gorm.DB, assignmentID int) (*models.Assignment, error) {
	var assignment models.Assignment
	if err := db.First(&assignment, assignmentID).Error; err != nil {
		return nil, err
	}
	return &assignment, nil
}

func CreateUserAssignment(db *gorm.DB, userAssignment *models.UserAssignment) error {
	return db.Create(userAssignment).Error
}
