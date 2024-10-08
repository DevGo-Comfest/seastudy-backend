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
	CourseID     int                `gorm:"type:int;not null"`
	IsLocked     *bool              `gorm:"-" json:"is_locked"`
	Materials    []SyllabusMaterial `gorm:"foreignKey:SyllabusID;constraint:OnDelete:CASCADE"`
	Assignments  []Assignment       `gorm:"foreignKey:SyllabusID;constraint:OnDelete:CASCADE"`
}