package controllers

import (
	"fmt"
	"net/http"
	"os"
	"sea-study/api/models"
	"sea-study/constants"
	"sea-study/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateCourse(c *gin.Context, db *gorm.DB) {
	var input models.CourseInput

	// userID Set by UserMiddleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidUserID})
		return
	}

	input.PrimaryAuthor = userUUID

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidInput})
		return
	}

	course, err := service.CreateCourse(db, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToCreateCourse})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course created successfully", "course": course})
}

func GetAllCourses(c *gin.Context, db *gorm.DB) {
	courses, err := service.GetAllCourses(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToRetrieveCourses})
		return
	}

	c.JSON(http.StatusOK, gin.H{"courses": courses})
}

func GetCourse(c *gin.Context, db *gorm.DB) {
	courseIDParam := c.Param("course_id")
	courseID, err := strconv.Atoi(courseIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidCourseID})
		return
	}

	course, err := service.GetCourse(db, courseID)
	if err != nil {
		if err.Error() == "course not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrCourseNotFound})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToRetrieveCourses})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"course": course})
}

func UpdateCourse(c *gin.Context, db *gorm.DB) {
	courseIDParam := c.Param("course_id")
	courseID, err := strconv.Atoi(courseIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidCourseID})
		return
	}

	var input models.CourseInput

	userID := c.GetString("userID")
	input.PrimaryAuthor = uuid.MustParse(userID)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidInput})
		return
	}

	course, err := service.UpdateCourse(db, courseID, &input)
	if err != nil {
		if err.Error() == "course not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrCourseNotFound})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToUpdateCourse})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course updated successfully", "course": course})
}

func DeleteCourse(c *gin.Context, db *gorm.DB) {
	courseIDParam := c.Param("course_id")
	courseID, err := strconv.Atoi(courseIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidCourseID})
		return
	}

	err = service.DeleteCourse(db, courseID)
	if err != nil {
		if err.Error() == "course not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrCourseNotFound})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToDeleteCourse})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}

func UploadCourseImage(c *gin.Context, db *gorm.DB) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrFailedToUploadImage})
		return
	}

	extension := file.Filename[len(file.Filename)-4:]

	imageID := uuid.New().String()

	filePath := fmt.Sprintf("uploads/%s%s", imageID, extension)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToSaveImage})
		return
	}

	hostURL := os.Getenv("HOST_URL")
	imageURL := fmt.Sprintf("%s/%s", hostURL, filePath)

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully", "image_url": imageURL})
}

func AddCourseInstructors(c *gin.Context, db *gorm.DB) {
	courseIDParam := c.Param("course_id")
	courseID, err := strconv.Atoi(courseIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidCourseID})
		return
	}

	// Get the user ID from the middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
		return
	}
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidUserID})
		return
	}

	// Check if the user is the creator of the course
	course, err := service.GetCourse(db, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToRetrieveCourses})
		return
	}

	if course.PrimaryAuthor != userUUID {
		c.JSON(http.StatusForbidden, gin.H{"error": constants.ErrUnauthorized})
		return
	}

	var instructorIDs models.InstructorIDs
	if err := c.ShouldBindJSON(&instructorIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidInput})
		return
	}

	if err := service.AddCourseInstructors(db, courseID, instructorIDs.InstructorIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToUpdateCourse})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Instructors added to course successfully"})
}

func SearchCourses(c *gin.Context, db *gorm.DB) {
	query := c.Query("q")                     // Search query
	category := c.Query("category")           // Course category filter
	difficulty := c.Query("difficulty_level") // Difficulty level filter
	ratingStr := c.Query("rating")            // Rating filter

	// Convert rating to integer
	var rating int
	var err error
	if ratingStr != "" {
		rating, err = strconv.Atoi(ratingStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidRating})
			return
		}
	}

	courses, err := service.SearchCourses(db, query, category, difficulty, rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToRetrieveCourses})
		return
	}

	c.JSON(http.StatusOK, gin.H{"courses": courses})
}

func GetMyCourse(c *gin.Context, db *gorm.DB) {
	userID := c.GetString("userID")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidUserID})
		return
	}

	course, err := service.GetCoursesByUser(db, userUUID)
	if err != nil {
		if err.Error() == "course not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrCourseNotFound})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrFailedToRetrieveCourses})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"course": course})
}

func ActivateCourse(c *gin.Context, db *gorm.DB) {
	courseID, err := strconv.Atoi(c.Param("course_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrInvalidCourseID})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
		return
	}

	course, err := service.ActivateCourse(db, userID.(string), courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course activated successfully", "course": course})
}
