package service

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"sea-study/api/models"
)

func CreateCourse(db *gorm.DB, input *models.CourseInput) (*models.Course, error) {
	categoryEnum := models.CategoryEnum(input.Category)
	course := &models.Course{
		Title:           input.Title,
		Description:     input.Description,
		Price:           input.Price,
		Category:        categoryEnum,
		ImageURL:        input.ImageURL,
		DifficultyLevel: input.DifficultyLevel,
		UserID:          input.UserID,
	}
	if err := db.Create(course).Error; err != nil {
		return nil, err
	}
	return course, nil
}

func GetAllCourses(db *gorm.DB) ([]models.Course, error) {
	var courses []models.Course
	result := db.Where("is_deleted = ? AND status = ?", false, models.ActiveStatus).Find(&courses)
	if result.Error != nil {
		return nil, result.Error
	}
	return courses, nil
}

func GetCourse(db *gorm.DB, courseID int) (*models.Course, error) {
	var course models.Course
	result := db.Preload("Syllabuses", func(db *gorm.DB) *gorm.DB {
		return db.Order("syllabuses.order")
	}).Where("course_id = ? AND is_deleted = ? AND status = ?", courseID, false, models.ActiveStatus).First(&course)
	if result.Error != nil {
		log.Println(result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("course not found")
		}
		return nil, result.Error
	}
	return &course, nil
}

func UpdateCourse(db *gorm.DB, courseID int, input *models.CourseInput) (*models.Course, error) {
	var course models.Course
	result := db.Where("course_id = ? AND is_deleted = ?", courseID, false).First(&course)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("course not found or has been deleted")
		}
		return nil, result.Error
	}
	categoryEnum := models.CategoryEnum(input.Category)
	course.Title = input.Title
	course.Description = input.Description
	course.Price = input.Price
	course.Category = categoryEnum
	course.ImageURL = input.ImageURL
	course.DifficultyLevel = input.DifficultyLevel
	if err := db.Save(&course).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func DeleteCourse(db *gorm.DB, courseID int) error {
	var course models.Course
	result := db.Where("course_id = ? AND is_deleted = ?", courseID, false).First(&course)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("course not found or already deleted")
		}
		return result.Error
	}

	// Soft delete by updating is_deleted to true
	if err := db.Model(&course).Update("is_deleted", true).Error; err != nil {
		return err
	}
	return nil
}

func AddCourseInstructors(db *gorm.DB, courseID int, instructorIDs []uuid.UUID) error {
	// First, check if the course exists and is not deleted
	var course models.Course
	if err := db.Where("course_id = ? AND is_deleted = ?", courseID, false).First(&course).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("course not found or has been deleted")
		}
		return err
	}

	for _, instructorID := range instructorIDs {
		courseInstructor := &models.CourseInstructor{
			CourseID:     courseID,
			InstructorID: instructorID,
		}
		if err := db.Create(courseInstructor).Error; err != nil {
			return err
		}
	}
	return nil
}

func SearchCourses(db *gorm.DB, query, category, difficulty string, rating int) ([]models.Course, error) {
	var courses []models.Course
	queryBuilder := db.Model(&models.Course{}).
		Where("is_deleted = ? AND status = ?", false, models.ActiveStatus)
	if query != "" {
		queryBuilder = queryBuilder.Where("title ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%")
	}
	if category != "" {
		queryBuilder = queryBuilder.Where("category = ?", category)
	}
	if difficulty != "" {
		queryBuilder = queryBuilder.Where("difficulty_level = ?", difficulty)
	}
	if rating > 0 {
		queryBuilder = queryBuilder.Where("rating >= ?", rating)
	}
	err := queryBuilder.Preload("Syllabuses").Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}
