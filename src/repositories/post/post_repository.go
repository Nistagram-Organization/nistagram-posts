package post

import (
	"github.com/Nistagram-Organization/nistagram-shared/src/datasources"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/post"
	"gorm.io/gorm"
)

type PostRepository interface {
	GetAll() []post.Post
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