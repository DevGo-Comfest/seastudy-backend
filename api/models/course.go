package models

import (
	"time"
)



type Course struct {
    CourseID       int             `gorm:"primaryKey;autoIncrement"`
    Title          string          `gorm:"type:varchar(255)"`
    Description    string          `gorm:"type:text"`
    Price          int             `gorm:"type:int"`
    Category       string          `gorm:"type:varchar(255)"`
    DifficultyLevel DifficultyEnum `gorm:"type:course_difficulty_enum"`
    CreatedDate    time.Time       `gorm:"type:timestamp"`
    UpdatedAt      time.Time       `gorm:"type:timestamp"`
    Rating         int             `gorm:"type:int"`
    Status         CourseStatusEnum `gorm:"type:course_status_enum"`
    Syllabuses     []Syllabus      `gorm:"foreignKey:CourseID"`
    Enrollments    []Enrollment    `gorm:"foreignKey:CourseID"`
    Progresses     []UserProgress  `gorm:"foreignKey:CourseID"`
    ForumPosts     []ForumPost     `gorm:"foreignKey:CourseID"`
    Reviews        []CourseReview  `gorm:"foreignKey:CourseID"`
}
