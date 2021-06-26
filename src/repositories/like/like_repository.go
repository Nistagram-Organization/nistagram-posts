package like

import (
	"fmt"
	"github.com/Nistagram-Organization/agent-shared/src/utils/rest_error"
	"github.com/Nistagram-Organization/nistagram-shared/src/datasources"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/like"
	"gorm.io/gorm"
)

type LikeRepository interface {
	Create(*like.Like) rest_error.RestErr
	GetByUserAndPost(uint, uint) (*like.Like, rest_error.RestErr)
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
	if err := l.db.Take(&likeEntity, likeEntity.UserID, likeEntity.PostID).Error; err != nil {
		fmt.Sprintln(err)
		return nil, rest_error.NewNotFoundError(fmt.Sprintf("Error when trying to get like"))
	}
	return &likeEntity, nil
}

func (l *likesRepository) Create(like *like.Like) rest_error.RestErr {
	if err := l.db.Create(like).Error; err != nil {
		return rest_error.NewInternalServerError("Error when trying to like a post", err)
	}
	return nil
}