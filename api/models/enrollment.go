package models

import (
	"time"

	"github.com/google/uuid"
)

type Enrollment struct {
    EnrollmentID int       `gorm:"primaryKey;autoIncrement"`
    UserID       uuid.UUID `gorm:"type:uuid;not null"`
    CourseID     int       `gorm:"type:int;not null"`
    DateEnrolled time.Time `gorm:"type:timestamp"`
}
