package service

import (
	"fmt"
	"log"
	"os"
	"sea-study/api/models"
	"sea-study/constants"
	"time"

	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

// Check if a user is already enrolled in a course
func IsUserEnrolled(db *gorm.DB, userID uuid.UUID, courseID int) (bool, error) {
	var enrollment models.Enrollment
	if err := db.Where("user_id = ? AND course_id = ?", userID, courseID).First(&enrollment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}


// Enroll a user in a course
func EnrollUser(db *gorm.DB, userID string, courseID int) (*models.Enrollment, error) {
    userUUID, err := uuid.Parse(userID)
    if err != nil {
        return nil, fmt.Errorf(constants.ErrInvalidUserID)
    }

    enrolled, err := IsUserEnrolled(db, userUUID, courseID)
    if err != nil {
        return nil, err
    }
    if enrolled {
        return nil, fmt.Errorf(constants.ErrUserAlreadyEnrolled)
    }

    var course models.Course
    if err := db.Where("course_id = ? AND status = ?", courseID, models.ActiveStatus).First(&course).Error; err != nil {
        return nil, err
    }

    balance, err := GetUserBalance(db, userUUID)
    if err != nil {
        return nil, err
    }
    if balance < float64(course.Price) {
        return nil, fmt.Errorf(constants.ErrInsufficientBalance)
    }

    // Start a transaction
    tx := db.Begin()

    if err := UpdateUserBalance(tx, userUUID, -float64(course.Price)); err != nil {
        tx.Rollback()
        return nil, err
    }

    enrollment := &models.Enrollment{
        UserID:       userUUID,
        CourseID:     courseID,
        DateEnrolled: time.Now(),
    }
    if err := tx.Create(enrollment).Error; err != nil {
        tx.Rollback()
        return nil, fmt.Errorf(constants.ErrFailedToCreateEnrollment)
    }

    tx.Commit()

    var username string
    if err := db.Table("users").Select("name").Where("user_id = ?", userUUID).Scan(&username).Error; err != nil {
        log.Printf("Failed to get username for user %s: %v", userUUID, err)
        username = "Student"
    }

    // Send enrollment notification emails to all instructors
    instructorEmails, err := getInstructorEmails(db, courseID)
    if err != nil {
        log.Printf("Failed to get instructor emails: %v", err)
    } else {
        for _, email := range instructorEmails {
            if err := sendEnrollmentEmail(email, course.Title, username); err != nil {
                log.Printf("Failed to send email to %s: %v", email, err)
            }
        }
    }


    return enrollment, nil
}

// Get instructor emails for a course
func getInstructorEmails(db *gorm.DB, courseID int) ([]string, error) {
    var instructorEmails []string

    // Get emails from CourseInstructor table
    err := db.Table("users").
        Select("users.email").
        Joins("JOIN course_instructors ON users.user_id = course_instructors.instructor_id").
        Where("course_instructors.course_id = ?", courseID).
        Scan(&instructorEmails).Error

    if err != nil {
        return nil, err
    }

    // Get email of the primary author
    var primaryAuthorEmail string
    err = db.Table("users").
        Select("email").
        Joins("JOIN courses ON users.user_id = courses.primary_author").
        Where("courses.course_id = ?", courseID).
        Scan(&primaryAuthorEmail).Error

    if err != nil {
        return nil, err
    }

    instructorEmails = append(instructorEmails, primaryAuthorEmail)

    return instructorEmails, nil
}

func sendEnrollmentEmail(email, courseTitle, username string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", os.Getenv("EMAIL_USERNAME"))
    m.SetHeader("To", email)
    m.SetHeader("Subject", "New Enrollment Notification")

    plainTextContent := fmt.Sprintf("Dear Instructor,\n\nUser %s has enrolled in your course: %s.\n\nBest regards,\nSea Study Team", username, courseTitle)
    htmlContent := fmt.Sprintf(`
        <p>Dear Instructor,</p>
        <p>User <strong>%s</strong> has enrolled in your course: <strong>%s</strong>.</p>
        <p>Best regards,<br>Sea Study Team</p>
    `, username, courseTitle)
    m.SetBody("text/plain", plainTextContent)
    m.AddAlternative("text/html", htmlContent)

    d := gomail.NewDialer(
        os.Getenv("SMTP_HOST"),
        587,
        os.Getenv("EMAIL_USERNAME"),
        os.Getenv("EMAIL_PASSWORD"),
    )

    err := d.DialAndSend(m)
    if err != nil {
        log.Printf("Failed to send email to %s: %v", email, err)
        return err
    }

    log.Printf("Email successfully sent to %s for course %s", email, courseTitle)
    return nil
}



func GetEnrolledCourses(db *gorm.DB, userID string) ([]models.Course, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidUserID)
	}

	var courses []models.Course
	err = db.Joins("JOIN enrollments ON enrollments.course_id = courses.course_id").
		Where("enrollments.user_id = ? AND courses.status = ?", userUUID, models.ActiveStatus).
		Order("enrollments.date_enrolled DESC").
		Find(&courses).Error
	if err != nil {
		return nil, fmt.Errorf(constants.ErrFailedToRetrieveEnrolledCourses)
	}

	return courses, nil
}
