package service

import (
	"errors"
	"fmt"
	"math"
	"sea-study/api/models"
	"sea-study/constants"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type StudentProgress struct {
	UserID           uuid.UUID `json:"user_id"`
	UserName         string    `json:"user_name"`
	ProgressPercentage int     `json:"progress_percentage"`
	LastAccessed     time.Time `json:"last_accessed"`
}


func UpdateUserProgress(db *gorm.DB, userID string, courseID, syllabusID int) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New(constants.ErrInvalidUserID)
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
			return errors.New(constants.ErrIncompletePreviousSyllabus)
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
		return 0, errors.New(constants.ErrInvalidUserID)
	}

	var totalSyllabuses int64
	if err := db.Model(&models.Syllabus{}).Where("course_id = ?", courseID).Count(&totalSyllabuses).Error; err != nil {
		return 0, err
	}

	if totalSyllabuses == 0 {
		return 0, errors.New(constants.ErrNoSyllabusesFound)
	}

	var completedSyllabuses int64
	if err := db.Model(&models.UserProgress{}).Where("user_id = ? AND course_id = ? AND status = ?", userUUID, courseID, models.Completed).Count(&completedSyllabuses).Error; err != nil {
		return 0, err
	}

	progressPercentage := int(math.Floor((float64(completedSyllabuses) / float64(totalSyllabuses)) * 100))

	return progressPercentage, nil
}


func GetStudentsProgressForCourse(db *gorm.DB, courseID int, instructorID string) ([]StudentProgress, error) {
	userUUID, err := uuid.Parse(instructorID)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidUserID)
	}

	// Check if the user is an instructor or the primary author of the course
	var count int64
	err = db.Table("course_instructors").
		Where("course_id = ? AND instructor_id = ?", courseID, userUUID).
		Count(&count).Error

	if err != nil {
		return nil, err
	}

	// If not found as an instructor, check if the user is the primary author
	if count == 0 {
		err = db.Model(&models.Course{}).
			Where("course_id = ? AND primary_author = ?", courseID, userUUID).
			Count(&count).Error

		if err != nil || count == 0 {
			return nil, fmt.Errorf(constants.ErrUnauthorized)
		}
	}

	// Get the total number of syllabuses for the course
	var totalSyllabuses int64
	if err := db.Model(&models.Syllabus{}).Where("course_id = ?", courseID).Count(&totalSyllabuses).Error; err != nil {
		return nil, fmt.Errorf(constants.ErrNoSyllabusesFound)
	}

	if totalSyllabuses == 0 {
		return nil, fmt.Errorf(constants.ErrNoSyllabusesFound)
	}

	// Retrieve the student progress
	var studentProgressList []StudentProgress
	err = db.Table("user_progress").
		Select(`user_progress.user_id, users.name as user_name, 
				COUNT(CASE WHEN user_progress.status = ? THEN 1 END) * 100 / ? as progress_percentage, 
				MAX(user_progress.last_accessed) as last_accessed`, 
				models.Completed, totalSyllabuses).
		Joins("JOIN users ON user_progress.user_id = users.user_id").
		Where("user_progress.course_id = ?", courseID).
		Group("user_progress.user_id, users.name").
		Scan(&studentProgressList).Error

	if err != nil {
		return nil, fmt.Errorf(constants.ErrFailedToRetrieveProgress)
	}

	return studentProgressList, nil
}

