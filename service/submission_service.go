package service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sea-study/api/models"
	"time"

	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

func CreateSubmission(db *gorm.DB, submission *models.Submission) (*models.Submission, error) {
    var assignmentTitle string
    err := db.Transaction(func(tx *gorm.DB) error {
        // Get the assignment based on assignment id in submission
        var assignment models.Assignment
        if err := tx.Where("assignment_id = ?", submission.AssignmentID).First(&assignment).Error; err != nil {
            return err
        }
        assignmentTitle = assignment.Title

        // Get the user assignment
        var userAssignment models.UserAssignment
        if err := tx.Where("assignment_id = ? AND user_id = ?", submission.AssignmentID, submission.UserID).First(&userAssignment).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                return fmt.Errorf("user assignment not found for assignment ID %d and user ID %s", submission.AssignmentID, submission.UserID)
            }
            return err
        }

        // Check if the submission is late
        submission.IsLate = time.Now().After(userAssignment.DueDate)

        // Get the syllabus based on the assignment's syllabus_id
        var syllabus models.Syllabus
        if err := tx.Where("syllabus_id = ?", assignment.SyllabusID).First(&syllabus).Error; err != nil {
            return err
        }

        // Create the submission
        if err := tx.Create(submission).Error; err != nil {
            return err
        }

        // Update or create the user progress
        var progress models.UserProgress
        result := tx.Where("user_id = ? AND course_id = ? AND syllabus_id = ?",
            submission.UserID, syllabus.CourseID, syllabus.SyllabusID).First(&progress)
        if result.Error != nil {
            if errors.Is(result.Error, gorm.ErrRecordNotFound) {
                // Create new progress if not found
                progress = models.UserProgress{
                    UserID:       submission.UserID,
                    CourseID:     syllabus.CourseID,
                    SyllabusID:   syllabus.SyllabusID,
                    Status:       models.ProgressStatusEnum(models.Completed),
                    LastAccessed: time.Now(),
                }
                if err := tx.Create(&progress).Error; err != nil {
                    return err
                }
            } else {
                return result.Error
            }
        } else {
            // Update existing progress
            if err := tx.Model(&progress).Updates(models.UserProgress{
                Status:       models.ProgressStatusEnum(models.Completed),
                LastAccessed: time.Now(),
            }).Error; err != nil {
                return err
            }
        }
        return nil
    })
    if err != nil {
        return nil, err
    }

    // Send submission notification emails to all instructors asynchronously
    go func() {
        courseID := getCourseIDBySyllabusID(db, submission.AssignmentID)
        if courseID != 0 {
            instructorEmails, err := getInstructorEmails(db, courseID)
            if err != nil {
                log.Printf("Failed to get instructor emails: %v", err)
                return
            }
            var studentName string
            if err := db.Table("users").Select("name").Where("user_id = ?", submission.UserID).Scan(&studentName).Error; err != nil {
                log.Printf("Failed to get student name for user %s: %v", submission.UserID, err)
                studentName = "Student"
            }
            for _, email := range instructorEmails {
                if err := sendSubmissionEmail(email, studentName, assignmentTitle); err != nil {
                    log.Printf("Failed to send email to %s: %v", email, err)
                }
            }
        }
    }()

    return submission, nil
}

// Get the course ID by syllabus ID
func getCourseIDBySyllabusID(db *gorm.DB, syllabusID int) int {
    var syllabus models.Syllabus
    if err := db.Where("syllabus_id = ?", syllabusID).First(&syllabus).Error; err != nil {
        log.Printf("Failed to get course ID by syllabus ID: %v", err)
        return 0
    }
    return syllabus.CourseID
}

func sendSubmissionEmail(email, studentName, assignmentTitle string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", os.Getenv("EMAIL_USERNAME"))
    m.SetHeader("To", email)
    m.SetHeader("Subject", "New Assignment Submission Notification")

    plainTextContent := fmt.Sprintf("Dear Instructor,\n\nStudent %s has submitted an assignment titled: %s.\n\nBest regards,\nSea Study Team", studentName, assignmentTitle)
    htmlContent := fmt.Sprintf(`
        <p>Dear Instructor,</p>
        <p>Student <strong>%s</strong> has submitted an assignment titled: <strong>%s</strong>.</p>
        <p>Best regards,<br>Sea Study Team</p>
    `, studentName, assignmentTitle)
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

    log.Printf("Email successfully sent to %s for assignment %s", email, assignmentTitle)
    return nil
}


func UpdateSubmission(db *gorm.DB, submission *models.Submission) (*models.Submission, error) {
	if err := db.Save(submission).Error; err != nil {
		return nil, err
	}
	return submission, nil
}

func DeleteSubmission(db *gorm.DB, submissionID int) error {
	result := db.Delete(&models.Submission{}, submissionID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("submission not found")
	}
	return nil
}

func GetSubmissionByID(db *gorm.DB, submissionID int) (*models.Submission, error) {
	var submission models.Submission
	if err := db.First(&submission, submissionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("submission not found")
		}
		return nil, err
	}
	return &submission, nil
}

func GetSubmissionByUserAndAssignment(db *gorm.DB, userID uuid.UUID, assignmentID int) (*models.Submission, error) {
	var submission models.Submission
	result := db.Where("user_id = ? AND assignment_id = ?", userID, assignmentID).First(&submission)
	if result.Error != nil {
		// No submission found
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &submission, nil
}
