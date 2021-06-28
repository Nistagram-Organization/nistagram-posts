package post

import (
	"fmt"
	"github.com/Nistagram-Organization/nistagram-shared/src/datasources"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"gorm.io/gorm"
)

type PostRepository interface {
	GetAll() []post.Post
	Get(uint) (*post.Post, rest_error.RestErr)
	Update(*post.Post) rest_error.RestErr
}

type postsRepository struct {
	db *gorm.DB
}

func NewPostRepository(databaseClient datasources.DatabaseClient) PostRepository {
	return &postsRepository{
		databaseClient.GetClient(),
	}
}

func (p *postsRepository) GetAll() []post.Post {
	var collection []post.Post
	if err := p.db.Find(&collection).Error; err != nil {
		return []post.Post{}
	}
	return collection
}

func (p *postsRepository) Get(id uint) (*post.Post, rest_error.RestErr) {
	post := post.Post{
		ID: id,
	}
	if err := p.db.Take(&post, post.ID).Error; err != nil {
		fmt.Sprintln(err)
		return nil, rest_error.NewNotFoundError(fmt.Sprintf("Error when trying to get post with id %d", post.ID))
	}
	return &post, nil
}

func (p *postsRepository) Update(post *post.Post) rest_error.RestErr {
	if err := p.db.Save(post).Error; err != nil {
		return rest_error.NewInternalServerError("Error when trying to update post", err)
	}
	return nil
}