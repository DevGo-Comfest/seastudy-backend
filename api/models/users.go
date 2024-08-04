package models

import (
	"time"

	"github.com/google/uuid"
)


type User struct {
    UserID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    Name      string    `gorm:"type:varchar(255)"`
    Email     string    `gorm:"type:varchar(255);unique"`
    Password  string    `gorm:"type:varchar(255)"`
    Balance   float64   `gorm:"type:decimal(10,2)"`
    Role      RoleEnum  `gorm:"type:role_enum"`
    CreatedAt time.Time `gorm:"type:timestamp"`
    UpdatedAt time.Time `gorm:"type:timestamp"`
    Topups    []TopupHistory `gorm:"foreignKey:UserID"`
    Courses   []Enrollment `gorm:"foreignKey:UserID"`
    Progresses []UserProgress `gorm:"foreignKey:UserID"`
    Assignments []Submission `gorm:"foreignKey:UserID"`
    Reviews    []CourseReview `gorm:"foreignKey:UserID"`
    ForumPosts []ForumPost `gorm:"foreignKey:UserID"`
}
