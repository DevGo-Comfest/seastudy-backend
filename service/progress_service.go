package service

import (
	"errors"
	"math"
	"sea-study/api/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UpdateUserProgress(db *gorm.DB, userID string, courseID, syllabusID int) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	var currentSyllabus models.Syllabus
	if err := db.Preload("Assignments").Where("syllabus_id = ?", syllabusID).First(&currentSyllabus).Error; err != nil {
		return err
	}

	var syllabuses []models.Syllabus
	if err := db.Where("course_id = ?", courseID).Order("\"order\"").Find(&syllabuses).Error; err != nil {
		return err
	}

	var userProgresses []models.UserProgress
	if err := db.Where("user_id = ? AND course_id = ?", userUUID, courseID).Find(&userProgresses).Error; err != nil {
		return err
	}

	userProgressMap := make(map[int]models.UserProgress)
	for _, progress := range userProgresses {
		userProgressMap[progress.SyllabusID] = progress
	}

	// Check if the user has completed each previous syllabus in order
	for _, syllabus := range syllabuses {
		if syllabus.Order >= currentSyllabus.Order {
			break
		}
		if progress, exists := userProgressMap[syllabus.SyllabusID]; !exists || progress.Status != models.Completed {
			return errors.New("complete all previous syllabuses to open this one")
		}
	}
	
	status := models.Completed
	if len(currentSyllabus.Assignments) > 0 {
		status = models.InProgress
	}

	progress := models.UserProgress{
		UserID:               userUUID,
		CourseID:             courseID,
		SyllabusID:           syllabusID,
		Status:               status,
		LastAccessed:         time.Now(),
	}

	if err := db.Where("user_id = ? AND course_id = ? AND syllabus_id = ?", userUUID, courseID, syllabusID).Assign(progress).FirstOrCreate(&progress).Error; err != nil {
		return err
	}

	return nil
}


func GetUserCourseProgress(db *gorm.DB, userID string, courseID int) (int, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, errors.New("invalid user ID")
	}

	var totalSyllabuses int64
	if err := db.Model(&models.Syllabus{}).Where("course_id = ?", courseID).Count(&totalSyllabuses).Error; err != nil {
		return 0, err
	}

	if totalSyllabuses == 0 {
		return 0, errors.New("no syllabuses found for this course")
	}

	// Count completed syllabuses for the user in the course
	var completedSyllabuses int64
	if err := db.Model(&models.UserProgress{}).Where("user_id = ? AND course_id = ? AND status = ?", userUUID, courseID, models.Completed).Count(&completedSyllabuses).Error; err != nil {
		return 0, err
	}

	// Calculate the completion percentage
	progressPercentage := int(math.Floor((float64(completedSyllabuses) / float64(totalSyllabuses)) * 100))

	return progressPercentage, nil
}