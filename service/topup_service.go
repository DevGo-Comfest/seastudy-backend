package service

import (
	"sea-study/api/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateTopup(db *gorm.DB, userID string, amount float64, paymentMethod string) (*models.TopupHistory, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	topup := &models.TopupHistory{
		UserID:        userUUID,
		Amount:        amount,
		Status:        "pending",
		PaymentMethod: paymentMethod,
		CreatedDate:   time.Now(),
	}

	// Start a transaction
	tx := db.Begin()
	if err := tx.Create(topup).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := UpdateUserBalance(tx, userUUID, amount); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return topup, nil
}
