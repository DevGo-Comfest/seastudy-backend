package service

import (
	"sea-study/api/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// CreateSyllabus
func TestCreateSyllabus(t *testing.T) {
	db, mock := setupTestDB(t)

	syllabus := &models.Syllabus{
		CourseID:     1,
		InstructorID: uuid.New(),
		Title:        "New Syllabus",
		Description:  "Syllabus Description",
	}

	mock.ExpectQuery(`SELECT count\(\*\) FROM "course_instructors"`).
		WithArgs(syllabus.CourseID, syllabus.InstructorID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(`SELECT COALESCE\(MAX\("order"\), 0\) FROM "syllabuses" WHERE course_id = \$1`).
		WithArgs(syllabus.CourseID).
		WillReturnRows(sqlmock.NewRows([]string{"max"}).AddRow(3))

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "syllabuses" (.+) VALUES (.+)`).
		WillReturnRows(sqlmock.NewRows([]string{"syllabus_id"}).AddRow(1))
	mock.ExpectCommit()

	err := CreateSyllabus(db, syllabus)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// UpdateSyllabus
func TestUpdateSyllabus(t *testing.T) {
	db, mock := setupTestDB(t)

	userUUID := uuid.New()
	syllabus := &models.Syllabus{
		SyllabusID:   1,
		CourseID:     1,
		InstructorID: userUUID,
		Title:        "Original Title",
		Description:  "Original Description",
	}

	userID := userUUID.String()

	mock.ExpectQuery(`SELECT .+ FROM "syllabuses".+`).
		WithArgs(syllabus.SyllabusID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"syllabus_id", "course_id", "instructor_id", "title", "description"}).
			AddRow(syllabus.SyllabusID, syllabus.CourseID, syllabus.InstructorID, syllabus.Title, syllabus.Description))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "syllabuses" SET`).
		WithArgs("Updated Description", "Updated Title", syllabus.SyllabusID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updatedSyllabus := &models.Syllabus{
		Title:       "Updated Title",
		Description: "Updated Description",
	}

	err := UpdateSyllabus(db, syllabus.SyllabusID, updatedSyllabus, userID)

	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updatedSyllabus.Title)
	assert.Equal(t, "Updated Description", updatedSyllabus.Description)
	assert.NoError(t, mock.ExpectationsWereMet())
}
