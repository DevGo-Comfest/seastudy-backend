package service

import (
	"fmt"
	"time"

	"sea-study/api/models"
	"sea-study/constants"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ForumPostResponse struct {
	ForumPostID int       `json:"forum_post_id"`
	CourseID    int       `json:"course_id"`
	Content     string    `json:"content"`
	UserID      string    `json:"user_id"`
	UserName    string    `json:"user_name"`
	UserRole    string    `json:"user_role"`
	DatePosted  time.Time `json:"date_posted"`
}

func CreateForumPost(db *gorm.DB, userID string, courseID int, content string) (*models.ForumPost, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidUserID)
	}

	forumPost := &models.ForumPost{
		CourseID:   courseID,
		UserID:     userUUID,
		Content:    content,
		DatePosted: time.Now(),
	}

	if err := db.Create(forumPost).Error; err != nil {
		return nil, fmt.Errorf(constants.ErrFailedToCreateForumPost)
	}

	return forumPost, nil
}

func GetForumPosts(db *gorm.DB, courseID int) ([]ForumPostResponse, error) {
	var forumPostResponses []ForumPostResponse

	err := db.Table("forum_posts").
		Select("forum_posts.forum_post_id, forum_posts.course_id, forum_posts.content, forum_posts.user_id, forum_posts.date_posted, users.name as user_name, users.role as user_role").
		Joins("left join users on forum_posts.user_id = users.user_id").
		Where("forum_posts.course_id = ?", courseID).
		Order("forum_posts.date_posted desc").
		Scan(&forumPostResponses).Error

	if err != nil {
		return nil, fmt.Errorf(constants.ErrFailedToRetrievePosts)
	}

	return forumPostResponses, nil
}
