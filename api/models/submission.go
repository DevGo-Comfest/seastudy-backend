package models

import (
	"time"

	"github.com/google/uuid"
)

type Submission struct {
	SubmissionID int                  `gorm:"primaryKey;autoIncrement"`
	Status       SubmissionStatusEnum `gorm:"type:submission_status_enum;default:submitted"`
	Grade        int                  `gorm:"type:int"`
	ContentURL   string               `gorm:"type:varchar(255)"`
	IsLate       bool                 `gorm:"type:boolean"`
	AssignmentID int                  `gorm:"type:int;not null"`
	UserID       uuid.UUID            `gorm:"type:uuid;not null"`
	CreatedAt    time.Time            `gorm:"type:autoCreateTime"`
	UpdatedAt    time.Time            `gorm:"type:autoUpdateTime"`
}
