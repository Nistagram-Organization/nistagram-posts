package post

import (
	"github.com/Nistagram-Organization/nistagram-posts/src/clients/media_grpc_client"
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
	UnlikePost(string, uint) rest_error.RestErr
	DislikePost(d *dtos.LikeDislikeRequestDTO) rest_error.RestErr
	UndislikePost(string, uint) rest_error.RestErr
	ReportInappropriateContent(uint) rest_error.RestErr
	PostComment(*modelComment.Comment) rest_error.RestErr
	CreatePost(*dtos.CreatePostDTO) rest_error.RestErr
	GetInappropriateContent() []dtos.InappropriateContentReportDTO
	DecideOnContent(uint, bool) rest_error.RestErr
}

type postsService struct {
	postsRepository    post.PostRepository
	likesRepository    like.LikeRepository
	dislikesRepository dislike.DislikeRepository
	commentsRepository comment.CommentRepository
	mediaGrpcClient    media_grpc_client.MediaGrpcClient
}

func NewPostService(postsRepository post.PostRepository, likesRepository like.LikeRepository, dislikesRepository dislike.DislikeRepository,
	commentsRepository comment.CommentRepository, mediaGrpcClient media_grpc_client.MediaGrpcClient) PostService {
	return &postsService{
		postsRepository:    postsRepository,
		likesRepository:    likesRepository,
		dislikesRepository: dislikesRepository,
		commentsRepository: commentsRepository,
		mediaGrpcClient:    mediaGrpcClient,
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

	if _, getLikeErr := s.likesRepository.GetByUserAndPost(likeRequest.UserEmail, likeRequest.PostID); getLikeErr == nil {
		return rest_error.NewBadRequestError("Post already liked")
	}

	if _, getDislikeErr := s.dislikesRepository.GetByUserAndPost(likeRequest.UserEmail, likeRequest.PostID); getDislikeErr == nil {
		return rest_error.NewBadRequestError("Post already disliked")
	}

	likeEntity := modelLike.Like{
		UserEmail: likeRequest.UserEmail,
		PostID:    likeRequest.PostID,
	}

	return s.likesRepository.Create(&likeEntity)
}

func (s *postsService) DislikePost(dislikeRequest *dtos.LikeDislikeRequestDTO) rest_error.RestErr {
	if err := s.checkIfPostExists(dislikeRequest.PostID); err != nil {
		return err
	}

	if _, getDislikeErr := s.dislikesRepository.GetByUserAndPost(dislikeRequest.UserEmail, dislikeRequest.PostID); getDislikeErr == nil {
		return rest_error.NewBadRequestError("Post already disliked")
	}

	if _, getLikeErr := s.likesRepository.GetByUserAndPost(dislikeRequest.UserEmail, dislikeRequest.PostID); getLikeErr == nil {
		return rest_error.NewBadRequestError("Post already liked")
	}

	dislikeEntity := modelDislike.Dislike{
		UserEmail: dislikeRequest.UserEmail,
		PostID:    dislikeRequest.PostID,
	}

	return s.dislikesRepository.Create(&dislikeEntity)
}

func (s *postsService) UnlikePost(userEmail string, postId uint) rest_error.RestErr {
	if err := s.checkIfPostExists(postId); err != nil {
		return err
	}

	if _, getLikeErr := s.likesRepository.GetByUserAndPost(userEmail, postId); getLikeErr != nil {
		return getLikeErr
	}

	likeEntity := modelLike.Like{
		UserEmail: userEmail,
		PostID:    postId,
	}

	return s.likesRepository.Delete(&likeEntity)
}

func (s *postsService) UndislikePost(userEmail string, postId uint) rest_error.RestErr {
	if err := s.checkIfPostExists(postId); err != nil {
		return err
	}

	if _, getDislikeErr := s.dislikesRepository.GetByUserAndPost(userEmail, postId); getDislikeErr != nil {
		return getDislikeErr
	}

	dislikeEntity := modelDislike.Dislike{
		UserEmail: userEmail,
		PostID:    postId,
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

func (s *postsService) CreatePost(postDTO *dtos.CreatePostDTO) rest_error.RestErr {
	saveMediaRequest := dtos.SaveMediaRequest{
		Image: postDTO.Image,
	}

	var mediaID *uint
	var err error

	if mediaID, err = s.mediaGrpcClient.SaveMedia(saveMediaRequest); err != nil {
		return rest_error.NewInternalServerError("user grpc client error when saving media", err)
	}

	postEntity := modelPost.Post{
		Description:           postDTO.Description,
		UserEmail:             postDTO.UserEmail,
		MarkedAsInappropriate: false,
		Date:                  time_utils.Now(),
		MediaID:               *mediaID,
	}

	return s.postsRepository.Create(&postEntity)
}

func (s *postsService) GetInappropriateContent() []dtos.InappropriateContentReportDTO {
	markedAsInappropriate := s.postsRepository.GetInappropriateContent()

	if len(markedAsInappropriate) == 0 {
		return []dtos.InappropriateContentReportDTO{}
	}

	var collection []dtos.InappropriateContentReportDTO
	for i := 0; i < len(markedAsInappropriate); i++ {
		media, _ := s.mediaGrpcClient.GetMedia(markedAsInappropriate[i].MediaID)

		inappropriateContentReport := dtos.InappropriateContentReportDTO{
			Description: markedAsInappropriate[i].Description,
			AuthorEmail: markedAsInappropriate[i].UserEmail,
			Image:       media,
			PostID:      markedAsInappropriate[i].ID,
		}
		collection = append(collection, inappropriateContentReport)
	}

	return collection
}

func (s *postsService) DecideOnContent(id uint, delete bool) rest_error.RestErr {
	post, err := s.postsRepository.Get(id)
	if err != nil {
		return err
	}

	if delete {
		if err := s.postsRepository.Delete(post); err != nil {
			return err
		}
	} else {
		post.MarkedAsInappropriate = false
		if err := s.postsRepository.Update(post); err != nil {
			return err
		}
	}

	return nil
}
