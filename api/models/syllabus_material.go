package models

type SyllabusMaterial struct {
	SyllabusMaterialID int    `gorm:"primaryKey;autoIncrement"`
	Order              int    `gorm:"type:int;not null"`
	Title              string `gorm:"type:varchar(255)"`
	Description        string `gorm:"type:text"`
	URLMaterial        string `gorm:"type:varchar(255)"`
	TimeNeeded         string `gorm:"type:varchar(255)"`
	SyllabusID         int    `gorm:"type:int;not null"`
}
