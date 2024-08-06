package controllers

import (
	"net/http"
	"sea-study/api/models"
	"sea-study/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateCourse(c *gin.Context, db* gorm.DB) {
	var input models.CourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	course, err := service.CreateCourse(db, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Course created successfully", "course": course})
}

func GetAllCourses(c *gin.Context, db *gorm.DB) {
    courses, err := service.GetAllCourses(db)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve courses"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"courses": courses})
}

func GetCourse(c *gin.Context, db *gorm.DB) {
    courseIDParam := c.Param("course_id")
    courseID, err := strconv.Atoi(courseIDParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
        return
    }
    
    course, err := service.GetCourse(db, courseID)
    if err != nil {
        if err.Error() == "course not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve course"})
        }
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"course": course})
}

func UpdateCourse(c *gin.Context, db *gorm.DB) {
    courseIDParam := c.Param("course_id")
    courseID, err := strconv.Atoi(courseIDParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
        return
    }

    var input models.CourseInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    course, err := service.UpdateCourse(db, courseID, &input)
    if err != nil {
        if err.Error() == "course not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Course updated successfully", "course": course})
}

func DeleteCourse(c *gin.Context, db *gorm.DB) {
    courseIDParam := c.Param("course_id")
    courseID, err := strconv.Atoi(courseIDParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
        return
    }

    err = service.DeleteCourse(db, courseID)
    if err != nil {
        if err.Error() == "course not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}