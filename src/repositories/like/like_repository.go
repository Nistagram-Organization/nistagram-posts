package like

import (
	"fmt"
	"github.com/Nistagram-Organization/nistagram-shared/src/datasources"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/like"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"gorm.io/gorm"
)

type LikeRepository interface {
	Create(*like.Like) rest_error.RestErr
	GetByUserAndPost(string, uint) (*like.Like, rest_error.RestErr)
	Delete(*like.Like) rest_error.RestErr
	GetNumberOfLikes(uint) (int64, rest_error.RestErr)
}

type likesRepository struct {
	db *gorm.DB
}

func NewLikeRepository(databaseClient datasources.DatabaseClient) LikeRepository {
	return &likesRepository{
		databaseClient.GetClient(),
	}
}

func (l *likesRepository) GetByUserAndPost(userEmail string, postId uint) (*like.Like, rest_error.RestErr) {
	likeEntity := like.Like{
		UserEmail: userEmail,
		PostID: postId,
	}
	if err := l.db.Where("user_email = ? AND post_id = ?", userEmail, postId).First(&likeEntity).Error; err != nil {
		return nil, rest_error.NewNotFoundError(fmt.Sprintf("Post has not been liked by user"))
	}
	return &likeEntity, nil
}

func (l *likesRepository) Create(like *like.Like) rest_error.RestErr {
	if err := l.db.Create(like).Error; err != nil {
		return rest_error.NewInternalServerError("Error when trying to like a post", err)
	}
	return nil
}

func (l *likesRepository) Delete(like *like.Like) rest_error.RestErr {
	if err := l.db.Where("user_email = ? AND post_id = ?", like.UserEmail, like.PostID).Delete(like).Error; err != nil {
		return rest_error.NewInternalServerError("Error when trying to unlike a post", err)
	}
	return nil
}

func (l *likesRepository) GetNumberOfLikes(postID uint) (int64, rest_error.RestErr) {
	var numberOfLikes int64
	if err := l.db.Where("post_id = ?", postID).Count(&numberOfLikes).Error; err != nil {
		return -1, rest_error.NewInternalServerError("Error when trying to get number of likes", err)
	}
	return numberOfLikes, nil
}
