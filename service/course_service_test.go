package service

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sea-study/api/models"
	"testing"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	return gormDB, mock
}

// CreateCourse
func TestCreateCourse(t *testing.T) {
	// Setup
	db, mock := setupTestDB(t)

	// Create a sample CourseInput
	primaryAuthor := uuid.New()
	input := &models.CourseInput{
		Title:           "Android Basic Course",
		Description:     "This is a test course",
		Price:           500000,
		Category:        models.CategoryEnum("Android"),
		ImageURL:        "http://example.com/image.jpg",
		DifficultyLevel: models.DifficultyEnum("intermediate"),
		PrimaryAuthor:   primaryAuthor,
	}

	// Expected SQL query
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "courses"`).
		WithArgs(
			primaryAuthor,
			input.Title,
			input.Description,
			input.Price,
			input.Category,
			input.ImageURL,
			input.DifficultyLevel,
			sqlmock.AnyArg(), // CreatedDate
			sqlmock.AnyArg(), // UpdatedAt
			0.0,              // Rating
			"inactive",       // Status
			false,            // IsDeleted
		).
		WillReturnRows(sqlmock.NewRows([]string{"course_id"}).AddRow(1))
	mock.ExpectCommit()

	// Execute the function
	course, err := CreateCourse(db, input)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, course)
	assert.Equal(t, input.Title, course.Title)
	assert.Equal(t, input.Description, course.Description)
	assert.Equal(t, input.Price, course.Price)
	assert.Equal(t, input.Category, course.Category)
	assert.Equal(t, input.ImageURL, course.ImageURL)
	assert.Equal(t, input.DifficultyLevel, course.DifficultyLevel)
	assert.Equal(t, input.PrimaryAuthor, course.PrimaryAuthor)
	assert.Equal(t, 1, course.CourseID)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

// GetAllCourses
func TestGetAllCourses(t *testing.T) {
	db, mock := setupTestDB(t)

	rows := sqlmock.NewRows([]string{"course_id", "title", "description", "price"}).
		AddRow(1, "Course 1", "Description 1", 1000000).
		AddRow(2, "Course 2", "Description 2", 2000000)

	mock.ExpectQuery(`SELECT \* FROM "courses" WHERE is_deleted = \$1 AND status = \$2`).
		WithArgs(false, models.ActiveStatus).
		WillReturnRows(rows)

	courses, err := GetAllCourses(db)

	assert.NoError(t, err)
	assert.Len(t, courses, 2)
	assert.Equal(t, "Course 1", courses[0].Title)
	assert.Equal(t, "Course 2", courses[1].Title)
}

// GetCourse
func TestGetCourse(t *testing.T) {
	db, mock := setupTestDB(t)

	courseID := 1

	rows := sqlmock.NewRows([]string{"course_id", "title", "description", "price"}).
		AddRow(courseID, "Test Course", "Test Description", 100)

	mock.ExpectQuery(`SELECT \* FROM "courses" WHERE "courses"."course_id" = \$1 ORDER BY "courses"."course_id" LIMIT \$2`).
		WithArgs(courseID, 1).
		WillReturnRows(rows)

	course, err := GetCourse(db, courseID)

	assert.NoError(t, err)
	assert.NotNil(t, course)
	assert.Equal(t, "Test Course", course.Title)
}

// GetCourseDetail
func TestGetCourseDetail(t *testing.T) {
	db, mock := setupTestDB(t)

	courseID := 1
	userID := uuid.New()
	primaryAuthorID := uuid.New()

	courseRows := sqlmock.NewRows([]string{"course_id", "title", "description", "price", "primary_author", "is_deleted", "status"}).
		AddRow(courseID, "Test Course", "Test Description", 100, primaryAuthorID, false, models.ActiveStatus)

	mock.ExpectQuery(`SELECT \* FROM "courses" WHERE course_id = \$1 AND is_deleted = \$2 AND status = \$3 ORDER BY "courses"."course_id" LIMIT \$4`).
		WithArgs(courseID, false, models.ActiveStatus, 1).
		WillReturnRows(courseRows)

	syllabusRows := sqlmock.NewRows([]string{"syllabus_id", "course_id", "title", "order"}).
		AddRow(1, courseID, "Syllabus 1", 1).
		AddRow(2, courseID, "Syllabus 2", 2)

	mock.ExpectQuery(`SELECT \* FROM "syllabuses" WHERE "syllabuses"."course_id" = \$1 ORDER BY syllabuses.order`).
		WithArgs(courseID).
		WillReturnRows(syllabusRows)

	// Mock primary author query
	authorRows := sqlmock.NewRows([]string{"user_id", "name"}).
		AddRow(primaryAuthorID, "John Doe")

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE user_id = \$1 ORDER BY "users"."user_id" LIMIT \$2`).
		WithArgs(primaryAuthorID, 1).
		WillReturnRows(authorRows)

	// Mock enrollment query
	enrollmentRows := sqlmock.NewRows([]string{"enrollment_id", "course_id", "user_id"}).
		AddRow(1, courseID, userID)

	mock.ExpectQuery(`SELECT \* FROM "enrollments" WHERE course_id = \$1 AND user_id = \$2 ORDER BY "enrollments"."enrollment_id" LIMIT \$3`).
		WithArgs(courseID, userID, 1).
		WillReturnRows(enrollmentRows)

	// Mock user progress query
	progressRows := sqlmock.NewRows([]string{"progress_id", "course_id", "user_id", "syllabus_id", "status"}).
		AddRow(1, courseID, userID, 1, models.Completed).
		AddRow(2, courseID, userID, 2, models.InProgress)

	mock.ExpectQuery(`SELECT \* FROM "user_progresses" WHERE course_id = \$1 AND user_id = \$2 ORDER BY syllabus_id`).
		WithArgs(courseID, userID).
		WillReturnRows(progressRows)

	course, err := GetCourseDetail(db, courseID, userID)

	assert.NoError(t, err)
	assert.NotNil(t, course)
	assert.Equal(t, "Test Course", course.Title)
	assert.Equal(t, "John Doe", course.PrimaryAuthorName)
	assert.Len(t, course.Syllabuses, 2)
	assert.False(t, *course.Syllabuses[0].IsLocked)
	assert.False(t, *course.Syllabuses[1].IsLocked)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// UpdateCourse
func TestUpdateCourse(t *testing.T) {
	db, mock := setupTestDB(t)

	courseID := 1
	input := &models.CourseInput{
		Title:           "Updated Course",
		Description:     "Updated Description",
		Price:           1000000,
		Category:        models.CategoryEnum("Android"),
		ImageURL:        "http://example.com/updated.jpg",
		DifficultyLevel: models.DifficultyEnum("intermediate"),
	}

	mock.ExpectQuery(`SELECT \* FROM "courses" WHERE course_id = \$1 AND is_deleted = \$2 ORDER BY "courses"."course_id" LIMIT \$3`).
		WithArgs(courseID, false, 1).
		WillReturnRows(sqlmock.NewRows([]string{"course_id"}).AddRow(courseID))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "courses" SET`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	course, err := UpdateCourse(db, courseID, input)

	assert.NoError(t, err)
	assert.NotNil(t, course)
	assert.Equal(t, input.Title, course.Title)
	assert.Equal(t, input.Description, course.Description)
	assert.Equal(t, input.Price, course.Price)
}

// DeleteCourse
func TestDeleteCourse(t *testing.T) {
	db, mock := setupTestDB(t)

	courseID := 1

	mock.ExpectQuery(`SELECT \* FROM "courses" WHERE course_id = \$1 AND is_deleted = \$2 ORDER BY "courses"."course_id" LIMIT \$3`).
		WithArgs(courseID, false, 1). // Correct the argument count to match the query
		WillReturnRows(sqlmock.NewRows([]string{"course_id"}).AddRow(courseID))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "courses" SET`).
		WithArgs(true, models.InactiveStatus, sqlmock.AnyArg(), courseID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := DeleteCourse(db, courseID)

	assert.NoError(t, err)
}

// AddCourseInstructors
func TestAddCourseInstructors(t *testing.T) {
	db, mock := setupTestDB(t)
	courseID := 1
	instructorIDs := []uuid.UUID{uuid.New(), uuid.New()} // Dynamic IDs

	mock.ExpectQuery(`SELECT \* FROM "courses" WHERE course_id = \$1 AND is_deleted = \$2 ORDER BY "courses"."course_id" LIMIT \$3`).
		WithArgs(courseID, false, 1).
		WillReturnRows(sqlmock.NewRows([]string{"course_id", "title"}).AddRow(courseID, "Test Course"))

	// Mock the queries for inserting instructors
	for _, instructorID := range instructorIDs {
		mock.ExpectBegin() // GORM wraps each Create in a transaction
		mock.ExpectQuery(`INSERT INTO "course_instructors" \("course_id","instructor_id"\) VALUES \(\$1,\$2\) RETURNING "course_instructor_id"`).
			WithArgs(courseID, instructorID).
			WillReturnRows(sqlmock.NewRows([]string{"course_instructor_id"}).AddRow(1))
		mock.ExpectCommit()
	}

	err := AddCourseInstructors(db, courseID, instructorIDs)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// SearchCourses
func TestSearchCourses(t *testing.T) {
	db, mock := setupTestDB(t)

	query := "test"
	category := "Android"
	difficulty := "intermediate"
	rating := 4

	courseRows := sqlmock.NewRows([]string{"course_id", "title"}).
		AddRow(1, "Test Course 1").
		AddRow(2, "Test Course 2")

	mock.ExpectQuery(`SELECT \* FROM "courses" WHERE \(is_deleted = \$1 AND status = \$2\) AND \(title ILIKE \$3 OR description ILIKE \$4\) AND category = \$5 AND difficulty_level = \$6 AND rating >= \$7`).
		WithArgs(false, "active", "%"+query+"%", "%"+query+"%", category, difficulty, rating).
		WillReturnRows(courseRows)

	syllabusRows := sqlmock.NewRows([]string{"id", "course_id", "title"}).
		AddRow(1, 1, "Syllabus 1").
		AddRow(2, 2, "Syllabus 2")

	mock.ExpectQuery(`SELECT \* FROM "syllabuses" WHERE "syllabuses"."course_id" IN \(\$1,\$2\)`).
		WithArgs(1, 2).
		WillReturnRows(syllabusRows)

	courses, err := SearchCourses(db, query, category, difficulty, rating)

	assert.NoError(t, err)
	assert.Len(t, courses, 2)
	assert.Equal(t, 1, courses[0].CourseID)
	assert.Equal(t, "Test Course 1", courses[0].Title)
	assert.Equal(t, 2, courses[1].CourseID)
	assert.Equal(t, "Test Course 2", courses[1].Title)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// GetCoursesByUser
func TestGetCoursesByUser(t *testing.T) {
	db, mock := setupTestDB(t)

	userID := uuid.New()

	rows := sqlmock.NewRows([]string{"course_id", "title"}).
		AddRow(1, "User Course 1").
		AddRow(2, "User Course 2")

	mock.ExpectQuery(`SELECT \* FROM "courses" WHERE primary_author = \$1 AND is_deleted = \$2`).
		WithArgs(userID, false).
		WillReturnRows(rows)

	courses, err := GetCoursesByUser(db, userID)

	assert.NoError(t, err)
	assert.Len(t, courses, 2)
}

// ActivateCourse
func TestActivateCourse(t *testing.T) {
	db, mock := setupTestDB(t)

	courseID := 1
	userID := uuid.New().String()

	mock.ExpectQuery(`SELECT \* FROM "courses" WHERE "courses"."course_id" = \$1 ORDER BY "courses"."course_id" LIMIT \$2`).
		WithArgs(courseID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"course_id", "primary_author"}).
			AddRow(courseID, userID))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "courses" SET`).
		WithArgs(models.ActiveStatus, sqlmock.AnyArg(), courseID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	course, err := ActivateCourse(db, userID, courseID)

	assert.NoError(t, err)
	assert.NotNil(t, course)
	assert.Equal(t, models.ActiveStatus, course.Status)
}

// GetPopularCourses
func TestGetPopularCourses(t *testing.T) {
	db, mock := setupTestDB(t)

	rows := sqlmock.NewRows([]string{"course_id", "title", "enrollment_count"}).
		AddRow(1, "Popular Course 1", 100).
		AddRow(2, "Popular Course 2", 80).
		AddRow(3, "Popular Course 3", 60).
		AddRow(4, "Popular Course 4", 40).
		AddRow(5, "Popular Course 5", 20)

	mock.ExpectQuery(`SELECT courses\.\*, COUNT\(enrollments\.enrollment_id\) as enrollment_count FROM "courses"`).
		WillReturnRows(rows)

	courses, err := GetPopularCourses(db)

	assert.NoError(t, err)
	assert.Len(t, courses, 5)
	assert.Equal(t, "Popular Course 1", courses[0].Title)
}
