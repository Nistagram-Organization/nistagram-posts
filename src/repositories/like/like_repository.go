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
	GetByUserAndPost(uint, uint) (*like.Like, rest_error.RestErr)
	Delete(*like.Like) rest_error.RestErr
}

type likesRepository struct {
	db *gorm.DB
}

func NewLikeRepository(databaseClient datasources.DatabaseClient) LikeRepository {
	return &likesRepository{
		databaseClient.GetClient(),
	}
}

func (l *likesRepository) GetByUserAndPost(userId uint, postId uint) (*like.Like, rest_error.RestErr) {
	likeEntity := like.Like{
		UserID: userId,
		PostID: postId,
	}
	if err := l.db.Where("user_id = ? AND post_id = ?", userId, postId).First(&likeEntity).Error; err != nil {
		fmt.Sprintln(err)
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
	if err :=l.db.Where("user_id = ? AND post_id = ?", like.UserID, like.PostID).Delete(like).Error; err != nil {
		return rest_error.NewInternalServerError("Error when trying to unlike a post", err)
	}
	return nil
}