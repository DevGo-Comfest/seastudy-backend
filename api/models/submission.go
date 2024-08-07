package models

import (
	"github.com/google/uuid"
)



type Submission struct {
    SubmissionID int              `gorm:"primaryKey;autoIncrement"`
    Status       SubmissionStatusEnum `gorm:"type:submission_status_enum"`
    Grade        int              `gorm:"type:int"`
    ContentURL   string           `gorm:"type:varchar(255)"`
    AssignmentID int              `gorm:"type:int;not null"`
    UserID       uuid.UUID        `gorm:"type:uuid;not null"`
}
