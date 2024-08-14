package service

import (
	"sea-study/api/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// CreateAssignment
func TestCreateAssignment(t *testing.T) {
	db, mock := setupTestDB(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO \"assignments"\ (.+) VALUES (.+)`).
		WillReturnRows(mock.NewRows([]string{"assignment_id"}).AddRow(1))
	mock.ExpectCommit()

	var assignment models.Assignment
	_, err := CreateAssignment(db, &assignment)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// OpenAssignment
func TestOpenAssignment(t *testing.T) {
	db, mock := setupTestDB(t)

	userID := uuid.New().String()
	assignmentID := 1

	assignmentRows := sqlmock.NewRows([]string{"assignment_id", "syllabus_id", "title", "description", "maximum_time"}).
		AddRow(assignmentID, 1, "Assignment 1", "Description 1", 7)

	mock.ExpectQuery(`SELECT \* FROM "assignments" WHERE "assignments"."assignment_id" = \$1 ORDER BY "assignments"."assignment_id" LIMIT \$2`).
		WithArgs(assignmentID, 1).
		WillReturnRows(assignmentRows)

	submissionRows := sqlmock.NewRows([]string{"submission_id", "assignment_id", "user_id"}).
		AddRow(1, assignmentID, uuid.New()).
		AddRow(2, assignmentID, uuid.New())

	mock.ExpectQuery(`SELECT \* FROM "submissions" WHERE "submissions"."assignment_id" = \$1`).
		WithArgs(1).
		WillReturnRows(submissionRows)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO \"user_assignments"\ (.+) VALUES (.+)`).
		WillReturnRows(mock.NewRows([]string{"user_assignment_id"}).AddRow(1))
	mock.ExpectCommit()

	err := OpenAssignment(db, userID, assignmentID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// GetAssignmentByID
func TestGetAssignmentByID(t *testing.T) {
	db, mock := setupTestDB(t)

	assignmentID := 1

	assignmentRows := sqlmock.NewRows([]string{"assignment_id", "syllabus_id", "title", "description", "maximum_time"}).
		AddRow(assignmentID, 1, "Assignment 1", "Description 1", 7)

	mock.ExpectQuery(`SELECT \* FROM "assignments" WHERE "assignments"."assignment_id" = \$1 ORDER BY "assignments"."assignment_id" LIMIT \$2`).
		WithArgs(assignmentID, 1).
		WillReturnRows(assignmentRows)

	submissionRows := sqlmock.NewRows([]string{"submission_id", "assignment_id", "user_id"}).
		AddRow(1, assignmentID, uuid.New()).
		AddRow(2, assignmentID, uuid.New())

	mock.ExpectQuery(`SELECT \* FROM "submissions" WHERE "submissions"."assignment_id" = \$1`).
		WithArgs(1).
		WillReturnRows(submissionRows)

	assignment, err := GetAssignmentByID(db, 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, "Assignment 1", assignment.Title)

}

// CreateUserAssignment
func TestCreateUserAssignment(t *testing.T) {
	db, mock := setupTestDB(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO \"user_assignments"\ (.+) VALUES (.+)`).
		WillReturnRows(mock.NewRows([]string{"user_assignment_id"}).AddRow(1))
	mock.ExpectCommit()

	var userAssignment models.UserAssignment
	err := CreateUserAssignment(db, &userAssignment)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// UpdateAssignment
func TestUpdateAssignment(t *testing.T) {
	db, mock := setupTestDB(t)

	assignment := &models.Assignment{
		AssignmentID: 1,
		SyllabusID:   1,
		Title:        "Original Title",
		Description:  "Original Description",
		MaximumTime:  7,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "assignments" SET`).
		WithArgs(
			assignment.SyllabusID,
			"Updated Title",
			"Updated Description",
			10,
			assignment.AssignmentID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	assignment.Title = "Updated Title"
	assignment.Description = "Updated Description"
	assignment.MaximumTime = 10

	updatedAssignment, err := UpdateAssignment(db, assignment)

	assert.NoError(t, err)
	assert.NotNil(t, updatedAssignment)
	assert.Equal(t, "Updated Title", updatedAssignment.Title)
	assert.Equal(t, "Updated Description", updatedAssignment.Description)
	assert.Equal(t, 10, updatedAssignment.MaximumTime)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// DeleteAssignment
func TestDeleteAssignment(t *testing.T) {
	db, mock := setupTestDB(t)

	assignmentID := 1

	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "assignments" WHERE "assignments"."assignment_id" = \$1`).
		WithArgs(assignmentID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := DeleteAssignment(db, assignmentID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
