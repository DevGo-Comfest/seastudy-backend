package service

import (
	"sea-study/constants"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// CreateForumPost
func TestCreateForumPost(t *testing.T) {
	db, mock := setupTestDB(t)

	userID := uuid.New()
	courseID := 1
	content := "Test forum post content"

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "forum_posts"`).
		WithArgs(
			courseID,
			userID,
			content,
			sqlmock.AnyArg(), // DatePosted
		).
		WillReturnRows(sqlmock.NewRows([]string{"forum_post_id"}).AddRow(1))
	mock.ExpectCommit()

	forumPost, err := CreateForumPost(db, userID.String(), courseID, content)

	assert.NoError(t, err)
	assert.NotNil(t, forumPost)
	assert.Equal(t, courseID, forumPost.CourseID)
	assert.Equal(t, userID, forumPost.UserID)
	assert.Equal(t, content, forumPost.Content)
	assert.NotZero(t, forumPost.DatePosted)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestCreateForumPostInvalidUserID(t *testing.T) {
	db, _ := setupTestDB(t)

	_, err := CreateForumPost(db, "invalid-uuid", 1, "content")

	assert.Error(t, err)
	assert.Equal(t, constants.ErrInvalidUserID, err.Error())
}

// GetForumPost
func TestGetForumPosts(t *testing.T) {
	db, mock := setupTestDB(t)

	courseID := 1
	userID := uuid.New()
	now := time.Now()

	rows := sqlmock.NewRows([]string{"forum_post_id", "course_id", "content", "user_id", "date_posted", "user_name", "user_role"}).
		AddRow(1, courseID, "Post 1", userID, now, "dimas", "user").
		AddRow(2, courseID, "Post 2", userID, now.Add(time.Hour), "hiya", "author")

	mock.ExpectQuery(`SELECT forum_posts.forum_post_id, forum_posts.course_id, forum_posts.content, forum_posts.user_id, forum_posts.date_posted, users.name as user_name, users.role as user_role FROM "forum_posts" left join users on forum_posts.user_id = users.user_id WHERE forum_posts.course_id = \$1 ORDER BY forum_posts.date_posted desc`).
		WithArgs(courseID).
		WillReturnRows(rows)

	forumPosts, err := GetForumPosts(db, courseID)

	assert.NoError(t, err)
	assert.Len(t, forumPosts, 2)
	assert.Equal(t, 1, forumPosts[0].ForumPostID)
	assert.Equal(t, courseID, forumPosts[0].CourseID)
	assert.Equal(t, "Post 1", forumPosts[0].Content)
	assert.Equal(t, userID.String(), forumPosts[0].UserID)
	assert.Equal(t, "dimas", forumPosts[0].UserName)
	assert.Equal(t, "user", forumPosts[0].UserRole)
	assert.Equal(t, now, forumPosts[0].DatePosted)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetForumPostsNoResults(t *testing.T) {
	db, mock := setupTestDB(t)

	courseID := 1

	mock.ExpectQuery(`SELECT forum_posts.forum_post_id, forum_posts.course_id, forum_posts.content, forum_posts.user_id, forum_posts.date_posted, users.name as user_name, users.role as user_role FROM "forum_posts" left join users on forum_posts.user_id = users.user_id WHERE forum_posts.course_id = \$1 ORDER BY forum_posts.date_posted desc`).
		WithArgs(courseID).
		WillReturnRows(sqlmock.NewRows([]string{"forum_post_id", "course_id", "content", "user_id", "date_posted", "user_name", "user_role"}))

	forumPosts, err := GetForumPosts(db, courseID)

	assert.NoError(t, err)
	assert.Len(t, forumPosts, 0)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}