package post

import (
	"github.com/Nistagram-Organization/nistagram-posts/src/dtos"
	"github.com/Nistagram-Organization/nistagram-posts/src/repositories/comment"
	"github.com/Nistagram-Organization/nistagram-posts/src/repositories/dislike"
	"github.com/Nistagram-Organization/nistagram-posts/src/repositories/like"
	"github.com/Nistagram-Organization/nistagram-posts/src/repositories/post"
	"github.com/Nistagram-Organization/nistagram-posts/src/time_utils"
	modelComment "github.com/Nistagram-Organization/nistagram-shared/src/model/comment"
	modelDislike "github.com/Nistagram-Organization/nistagram-shared/src/model/dislike"
	modelLike "github.com/Nistagram-Organization/nistagram-shared/src/model/like"
	modelPost "github.com/Nistagram-Organization/nistagram-shared/src/model/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
)

type PostService interface {
	GetAll() []modelPost.Post
	LikePost(*dtos.LikeDislikeRequestDTO) rest_error.RestErr
	UnlikePost(uint, uint) rest_error.RestErr
	DislikePost(d *dtos.LikeDislikeRequestDTO) rest_error.RestErr
	UndislikePost(uint, uint) rest_error.RestErr
	ReportInappropriateContent(uint) rest_error.RestErr
	PostComment(*modelComment.Comment) rest_error.RestErr
}

type postsService struct {
	postsRepository    post.PostRepository
	likesRepository    like.LikeRepository
	dislikesRepository dislike.DislikeRepository
	commentsRepository comment.CommentRepository
}

func NewPostService(postsRepository post.PostRepository, likesRepository like.LikeRepository, dislikesRepository dislike.DislikeRepository,
	commentsRepository comment.CommentRepository) PostService {
	return &postsService{
		postsRepository:    postsRepository,
		likesRepository:    likesRepository,
		dislikesRepository: dislikesRepository,
		commentsRepository: commentsRepository,
	}
}

func (s *postsService) checkIfPostExists(postId uint) rest_error.RestErr {
	_, err := s.postsRepository.Get(postId)
	if err != nil {
		return err
	}

	return nil
}

func (s *postsService) GetAll() []modelPost.Post {
	return s.postsRepository.GetAll()
}

func (s *postsService) LikePost(likeRequest *dtos.LikeDislikeRequestDTO) rest_error.RestErr {
	if err := s.checkIfPostExists(likeRequest.PostID); err != nil {
		return err
	}

	if _, getLikeErr := s.likesRepository.GetByUserAndPost(likeRequest.UserID, likeRequest.PostID); getLikeErr == nil {
		return rest_error.NewBadRequestError("Post already liked")
	}

	if _, getDislikeErr := s.dislikesRepository.GetByUserAndPost(likeRequest.UserID, likeRequest.PostID); getDislikeErr == nil {
		return rest_error.NewBadRequestError("Post already disliked")
	}

	likeEntity := modelLike.Like{
		UserID: likeRequest.UserID,
		PostID: likeRequest.PostID,
	}

	return s.likesRepository.Create(&likeEntity)
}

func (s *postsService) DislikePost(dislikeRequest *dtos.LikeDislikeRequestDTO) rest_error.RestErr {
	if err := s.checkIfPostExists(dislikeRequest.PostID); err != nil {
		return err
	}

	if _, getDislikeErr := s.dislikesRepository.GetByUserAndPost(dislikeRequest.UserID, dislikeRequest.PostID); getDislikeErr == nil {
		return rest_error.NewBadRequestError("Post already disliked")
	}

	if _, getLikeErr := s.likesRepository.GetByUserAndPost(dislikeRequest.UserID, dislikeRequest.PostID); getLikeErr == nil {
		return rest_error.NewBadRequestError("Post already liked")
	}

	dislikeEntity := modelDislike.Dislike{
		UserID: dislikeRequest.UserID,
		PostID: dislikeRequest.PostID,
	}

	return s.dislikesRepository.Create(&dislikeEntity)
}

func (s *postsService) UnlikePost(userId uint, postId uint) rest_error.RestErr {
	if err := s.checkIfPostExists(postId); err != nil {
		return err
	}

	if _, getLikeErr := s.likesRepository.GetByUserAndPost(userId, postId); getLikeErr != nil {
		return getLikeErr
	}

	likeEntity := modelLike.Like{
		UserID: userId,
		PostID: postId,
	}

	return s.likesRepository.Delete(&likeEntity)
}

func (s *postsService) UndislikePost(userId uint, postId uint) rest_error.RestErr {
	if err := s.checkIfPostExists(postId); err != nil {
		return err
	}

	if _, getDislikeErr := s.dislikesRepository.GetByUserAndPost(userId, postId); getDislikeErr != nil {
		return getDislikeErr
	}

	dislikeEntity := modelDislike.Dislike{
		UserID: userId,
		PostID: postId,
	}

	return s.dislikesRepository.Delete(&dislikeEntity)
}

func (s *postsService) ReportInappropriateContent(postId uint) rest_error.RestErr {
	postEntity, err := s.postsRepository.Get(postId)
	if err != nil {
		return err
	}

	if !postEntity.MarkedAsInappropriate {
		postEntity.MarkedAsInappropriate = true
		return s.postsRepository.Update(postEntity)
	} else {
		return nil
	}
}

func (s *postsService) PostComment(commentEntity *modelComment.Comment) rest_error.RestErr {
	if err := s.checkIfPostExists(commentEntity.PostID); err != nil {
		return err
	}
	commentEntity.Date = time_utils.Now()

	return s.commentsRepository.Create(commentEntity)
}
