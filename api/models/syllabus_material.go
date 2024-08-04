package models

type SyllabusMaterial struct {
	SyllabusMaterialID int    `gorm:"primaryKey;autoIncrement"`
	URLMaterial        string `gorm:"type:varchar(255)"`
	TimeNeeded         string `gorm:"type:varchar(255)"`
	SyllabusID         int    `gorm:"type:int;not null"`
}
