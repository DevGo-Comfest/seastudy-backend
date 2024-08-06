package models

import (
	"github.com/google/uuid"
)

type Syllabus struct {
	SyllabusID   int                `gorm:"primaryKey;autoIncrement"`
	Order        int                `gorm:"type:int;not null"`
	Title        string             `gorm:"type:varchar(255)"`
	Description  string             `gorm:"type:text"`
	InstructorID uuid.UUID          `gorm:"type:uuid;not null"`
	AssignmentID string             `gorm:"type:varchar(255)"`
	CourseID     int                `gorm:"type:int;not null"`
	Materials    []SyllabusMaterial `gorm:"foreignKey:SyllabusID"`
	Assignments  []Assignment       `gorm:"foreignKey:SyllabusID"`
}
