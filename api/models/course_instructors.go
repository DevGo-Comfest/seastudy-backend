package models

import "github.com/google/uuid"

type CourseInstructor struct {
	CourseInstructorID int     `gorm:"primaryKey;autoIncrement"`
	CourseID           int     `gorm:"type:int;not null"`
	InstructorID       uuid.UUID `gorm:"type:uuid;not null"`
}
