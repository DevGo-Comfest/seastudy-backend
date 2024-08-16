package service

import (
	"fmt"
	"sea-study/api/models"
	"sea-study/constants"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateTopup(db *gorm.DB, userID string, amount float64, paymentMethod string) (*models.TopupHistory, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidUserID)
	}

	topup := &models.TopupHistory{
		UserID:        userUUID,
		Amount:        amount,
		Status:        "completed",
		PaymentMethod: paymentMethod,
		CreatedDate:   time.Now(),
	}

	// Start a transaction
	tx := db.Begin()
	if err := tx.Create(topup).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf(constants.ErrFailedToCreateTopup)
	}

	if err := UpdateUserBalance(tx, userUUID, amount); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf(constants.ErrFailedToUpdateUserBalance)
	}

	tx.Commit()

	return topup, nil
}

func GetTopupHistory(db *gorm.DB, userID string) ([]models.TopupHistory, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidUserID)
	}

	var history []models.TopupHistory
	err = db.Where("user_id = ?", userUUID).Order("created_date desc").Find(&history).Error
	if err != nil {
		return nil, err
	}

	return history, nil
}