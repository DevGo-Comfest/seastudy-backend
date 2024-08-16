package models

import (
	"time"

	"github.com/google/uuid"
)

type UserProgress struct {
    UserProgressID     int              `gorm:"primaryKey;autoIncrement"`
    UserID             uuid.UUID        `gorm:"type:uuid;not null"`
    CourseID           int              `gorm:"type:int;not null"`
    SyllabusID         int              `gorm:"type:int;not null"`
    Status             ProgressStatusEnum `gorm:"type:progress_status_enum"`
    LastAccessed       time.Time        `gorm:"type:timestamp"`
}
