package models

import (
	"time"

	"github.com/google/uuid"
)

type CourseReview struct {
    CourseReviewID int       `gorm:"primaryKey;autoIncrement"`
    CourseID       int       `gorm:"type:int;not null"`
    FeedbackText   string    `gorm:"type:varchar(255)"`
    UserID         uuid.UUID `gorm:"type:uuid;not null"`
    CreatedAt      time.Time `gorm:"type:timestamp"`
    UpdatedAt      time.Time `gorm:"type:timestamp"`
    Rate           int       `gorm:"type:int"`
}
