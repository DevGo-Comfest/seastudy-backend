package models

import (
	"time"

	"github.com/google/uuid"
)

type ForumPost struct {
    ForumPostID int       `gorm:"primaryKey;autoIncrement"`
    CourseID    int       `gorm:"type:int;not null"`
    UserID      uuid.UUID `gorm:"type:uuid;not null"`
    Content     string    `gorm:"type:text"`
    DatePosted  time.Time `gorm:"type:timestamp"`
}
