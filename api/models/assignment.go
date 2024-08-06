package models

type Assignment struct {
	AssignmentID int          `gorm:"primaryKey;autoIncrement"`
	SyllabusID   int          `gorm:"type:int;not null"`
	Title        string       `gorm:"type:varchar(255)"`
	Description  string       `gorm:"type:text"`
	MaximumTime  int          `gorm:"type:int"`
	Submissions  []Submission `gorm:"foreignKey:AssignmentID;constraint:OnDelete:CASCADE"`
}