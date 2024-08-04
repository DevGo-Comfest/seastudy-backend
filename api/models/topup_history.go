package models

import (
	"time"

	"github.com/google/uuid"
)


type TopupHistory struct {
    TopupID       int            `gorm:"primaryKey;autoIncrement"`
    UserID        uuid.UUID      `gorm:"type:uuid;not null"`
    Amount        float64        `gorm:"type:decimal(10,2)"`
    Status        TopupStatusEnum `gorm:"type:topup_status_enum;default:'pending'"`
    PaymentMethod string         `gorm:"type:varchar(255)"`
    CreatedDate   time.Time      `gorm:"type:timestamp"`
}
