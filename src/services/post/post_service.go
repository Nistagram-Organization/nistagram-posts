package post

import (
	"github.com/Nistagram-Organization/nistagram-posts/src/dtos"
	"github.com/Nistagram-Organization/nistagram-posts/src/repositories/like"
	"github.com/Nistagram-Organization/nistagram-posts/src/repositories/post"
	modelLike "github.com/Nistagram-Organization/nistagram-shared/src/model/like"
	modelPost "github.com/Nistagram-Organization/nistagram-shared/src/model/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
)

type PostService interface {
	GetAll() []modelPost.Post
	LikePost(* dtos.LikeRequestDTO) rest_error.RestErr
	UnlikePost(*dtos.LikeRequestDTO) rest_error.RestErr
}

type postsService struct {
	postsRepository post.PostRepository
	likesRepository like.LikeRepository
}

func NewPostService(postsRepository post.PostRepository, likesRepository like.LikeRepository) PostService {
	return &postsService{
		postsRepository: postsRepository,
		likesRepository: likesRepository,
	}
}

func (s *postsService) GetAll() []modelPost.Post {
	return s.postsRepository.GetAll()
}

func (s *postsService) LikePost(likeRequest *dtos.LikeRequestDTO) rest_error.RestErr {
	if _, getLikeErr := s.likesRepository.GetByUserAndPost(likeRequest.UserID, likeRequest.PostID); getLikeErr == nil {
		return rest_error.NewBadRequestError("Post already liked")
	}

	likeEntity := modelLike.Like{
		UserID: likeRequest.UserID,
		PostID: likeRequest.PostID,
	}

	return s.likesRepository.Create(&likeEntity)
}

func (s *postsService) UnlikePost(likeRequest *dtos.LikeRequestDTO) rest_error.RestErr {
	if _, getLikeErr := s.likesRepository.GetByUserAndPost(likeRequest.UserID, likeRequest.PostID); getLikeErr != nil {
		return getLikeErr
	}

	likeEntity := modelLike.Like{
		UserID: likeRequest.UserID,
		PostID: likeRequest.PostID,
	}

	return s.likesRepository.Delete(&likeEntity)
}