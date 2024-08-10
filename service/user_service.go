package service

import (
	"fmt"
	"sea-study/api/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *models.User) error {
	return db.Create(user).Error
}

func GetUserByEmail(db *gorm.DB, user *models.User, email string) error {
	return db.Where("email = ?", email).First(user).Error
}

func GetUserByID(db *gorm.DB, user *models.User, userID uuid.UUID) error {
	return db.Where("user_id = ?", userID).First(user).Error
}

func ValidateUserExists(db *gorm.DB, userID uuid.UUID) error {
	var user models.User
	if err := db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserBalance(db *gorm.DB, userID uuid.UUID) (float64, error) {
	var user models.User
	if err := db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return 0, err
	}
	return user.Balance, nil
}

func UpdateUserBalance(db *gorm.DB, userID uuid.UUID, amount float64) error {
	var user models.User
	if err := db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	if amount < 0 && user.Balance < -amount {
		return fmt.Errorf("insufficient balance")
	}

	user.Balance += amount

	return db.Save(&user).Error
}
