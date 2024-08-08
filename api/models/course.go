package models

import (
	"time"

	"github.com/google/uuid"
)

type InstructorIDs struct {
	InstructorIDs []uuid.UUID `json:"instructor_ids"`
}

type CourseInput struct {
	Title           string         `json:"title" binding:"required"`
	Description     string         `json:"description" binding:"required"`
	Price           int            `json:"price" binding:"required"`
	Category        string         `json:"category" binding:"required"`
	ImageURL        string         `json:"image_url" binding:"required"`
	DifficultyLevel DifficultyEnum `json:"difficulty_level" binding:"required"`
	UserID          uuid.UUID      `json:"user_id" binding:"required"`
}

type Course struct {
	CourseID        int              `gorm:"primaryKey;autoIncrement"`
	UserID          uuid.UUID        `gorm:"type:uuid;not null"`
	Title           string           `gorm:"type:varchar(255)"`
	Description     string           `gorm:"type:text"`
	Price           int              `gorm:"type:int"`
	Category        string           `gorm:"type:varchar(255)"`
	ImageURL        string           `gorm:"type:text"`
	DifficultyLevel DifficultyEnum   `gorm:"type:course_difficulty_enum"`
	CreatedDate     time.Time        `gorm:"type:timestamp"`
	UpdatedAt       time.Time        `gorm:"type:timestamp"`
	Rating          int              `gorm:"type:int"`
	Status          CourseStatusEnum `gorm:"type:course_status_enum;default:inactive"`
	IsDeleted       bool             `gorm:"type:boolean;default:false"`
	Syllabuses      []Syllabus       `gorm:"foreignKey:CourseID;"`
	Enrollments     []Enrollment     `gorm:"foreignKey:CourseID"`
	Progresses      []UserProgress   `gorm:"foreignKey:CourseID"`
	ForumPosts      []ForumPost      `gorm:"foreignKey:CourseID;"`
	Reviews         []CourseReview   `gorm:"foreignKey:CourseID;"`
}
