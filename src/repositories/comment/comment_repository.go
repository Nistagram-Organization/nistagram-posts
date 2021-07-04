package comment

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/datasources"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/comment"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(*comment.Comment) rest_error.RestErr
	GetComments(uint) ([]comment.Comment, rest_error.RestErr)
}

type commentsRepository struct {
	db *gorm.DB
}

func NewCommentRepository(databaseClient datasources.DatabaseClient) CommentRepository {
	return &commentsRepository{
		databaseClient.GetClient(),
	}
}

func (c *commentsRepository) Create(comment *comment.Comment) rest_error.RestErr {
	if err := c.db.Create(comment).Error; err != nil {
		return rest_error.NewInternalServerError("Error when trying to post a comment", err)
	}
	return nil
}

func (c *commentsRepository) GetComments(postID uint) ([]comment.Comment, rest_error.RestErr) {
	var collection []comment.Comment

	if err := c.db.Where("post_id = ?", postID).Find(&collection).Error; err != nil {
		return nil, rest_error.NewInternalServerError("Error when trying to get post's comments", err)
	}

	return collection, nil
}
