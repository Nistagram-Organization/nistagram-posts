package post

import (
	"github.com/Nistagram-Organization/nistagram-posts/src/repositories/post"
	model "github.com/Nistagram-Organization/nistagram-shared/src/model/post"
)

type PostService interface {
	GetAll() []model.Post
}

type postsService struct {
	postsRepository post.PostRepository
}

func NewPostService(postsRepository post.PostRepository) PostService {
	return &postsService{
		postsRepository: postsRepository,
	}
}

func (s *postsService) GetAll() []model.Post {
	return s.postsRepository.GetAll()
}
