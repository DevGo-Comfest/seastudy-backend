package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID          uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name            string           `gorm:"type:varchar(255)" json:"name" binding:"required,min=2,max=255"`
	Email           string           `gorm:"type:varchar(255);unique" json:"email" binding:"required,email"`
	Password        string           `gorm:"type:varchar(255)" json:"password" binding:"required,min=8"`
	Balance         float64          `gorm:"type:decimal(10,2)"`
	Role            RoleEnum         `gorm:"type:role_enum;default:user"`
	CreatedAt       time.Time        `gorm:"type:timestamp"`
	UpdatedAt       time.Time        `gorm:"type:timestamp"`
	Topups          []TopupHistory   `gorm:"foreignKey:UserID"`
	Courses         []Enrollment     `gorm:"foreignKey:UserID"`
	Progresses      []UserProgress   `gorm:"foreignKey:UserID"`
	Submissions     []Submission     `gorm:"foreignKey:UserID"`
	UserAssignments []UserAssignment `gorm:"foreignKey:UserID"`
	Reviews         []CourseReview   `gorm:"foreignKey:UserID"`
	ForumPosts      []ForumPost      `gorm:"foreignKey:UserID"`
}
