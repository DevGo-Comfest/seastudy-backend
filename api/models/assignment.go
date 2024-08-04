package models

import (
	"time"
)

type Assignment struct {
    AssignmentID int       `gorm:"primaryKey;autoIncrement"`
    SyllabusID   int       `gorm:"type:int;not null"`
    Title        string    `gorm:"type:varchar(255)"`
    Description  string    `gorm:"type:text"`
    DueDate      time.Time `gorm:"type:timestamp"`
    Submissions  []Submission `gorm:"foreignKey:AssignmentID"`
}
