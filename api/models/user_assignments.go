package models

import (
	"time"

	"github.com/google/uuid"
)

type UserAssignment struct {
	UserAssignmentID int       `gorm:"primaryKey;autoIncrement"`
	AssignmentID     int       `gorm:"type:int;not null"`
	UserID           uuid.UUID `gorm:"type:uuid;not null"`
	DueDate          time.Time `gorm:"type:date"`
	CreatedAt        time.Time `gorm:"type:timestamp"`
}