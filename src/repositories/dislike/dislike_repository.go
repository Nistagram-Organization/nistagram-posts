package dislike

import (
	"fmt"
	"github.com/Nistagram-Organization/nistagram-shared/src/datasources"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/dislike"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"gorm.io/gorm"
)

type DislikeRepository interface {
	Create(*dislike.Dislike) rest_error.RestErr
	GetByUserAndPost(uint, uint) (*dislike.Dislike, rest_error.RestErr)
	Delete(*dislike.Dislike) rest_error.RestErr
}

type dislikesRepository struct {
	db *gorm.DB
}

func NewDislikeRepository(databaseClient datasources.DatabaseClient) DislikeRepository {
	return &dislikesRepository{
		databaseClient.GetClient(),
	}
}

func (d *dislikesRepository) GetByUserAndPost(userId uint, postId uint) (*dislike.Dislike, rest_error.RestErr) {
	dislikeEntity := dislike.Dislike{
		UserID: userId,
		PostID: postId,
	}
	if err := d.db.Where("user_id = ? AND post_id = ?", userId, postId).First(&dislikeEntity).Error; err != nil {
		return nil, rest_error.NewNotFoundError(fmt.Sprintf("Post has not been disliked by user"))
	}
	return &dislikeEntity, nil
}

func (d *dislikesRepository) Create(dislike *dislike.Dislike) rest_error.RestErr {
	if err := d.db.Create(dislike).Error; err != nil {
		return rest_error.NewInternalServerError("Error when trying to dislike a post", err)
	}
	return nil
}

func (d *dislikesRepository) Delete(dislike *dislike.Dislike) rest_error.RestErr {
	if err := d.db.Where("user_id = ? AND post_id = ?", dislike.UserID, dislike.PostID).Delete(dislike).Error; err != nil {
		return rest_error.NewInternalServerError("Error when trying to undislike a post", err)
	}
	return nil
}
