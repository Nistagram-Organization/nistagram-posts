package post

import (
	"github.com/Nistagram-Organization/nistagram-posts/src/clients/media_grpc_client"
	"github.com/Nistagram-Organization/nistagram-posts/src/clients/user_grpc_client"
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
	"regexp"
	"strings"
	"time"
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
	GetUsersPosts(string, string) ([]dtos.PostDTO, rest_error.RestErr)
}

type postsService struct {
	postsRepository    post.PostRepository
	likesRepository    like.LikeRepository
	dislikesRepository dislike.DislikeRepository
	commentsRepository comment.CommentRepository
	mediaGrpcClient    media_grpc_client.MediaGrpcClient
	userGrpcClient     user_grpc_client.UserGrpcClient
}

func NewPostService(postsRepository post.PostRepository, likesRepository like.LikeRepository, dislikesRepository dislike.DislikeRepository,
	commentsRepository comment.CommentRepository, mediaGrpcClient media_grpc_client.MediaGrpcClient, userGrpcClient user_grpc_client.UserGrpcClient) PostService {
	return &postsService{
		postsRepository:    postsRepository,
		likesRepository:    likesRepository,
		dislikesRepository: dislikesRepository,
		commentsRepository: commentsRepository,
		mediaGrpcClient:    mediaGrpcClient,
		userGrpcClient: userGrpcClient,
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

func (s *postsService) GetUsersPosts(userEmail string, loggedInUserEmail string) ([]dtos.PostDTO, rest_error.RestErr) {
	var postsDTOs []dtos.PostDTO
	var posts []modelPost.Post
	var postErr rest_error.RestErr

	// Get all users posts
	if posts, postErr = s.postsRepository.GetUsersPosts(userEmail); postErr != nil {
		return nil, postErr
	}

	layout := "02.01.2006. 03:04"
	for _, postEntity := range posts {
		description := s.ProcessTags(postEntity.Description)
		// Convert time to format dd.MM.yyyy. HH:mm
		t := time.Unix(postEntity.Date, 0)
		date := t.Format(layout)

		// GRPC call media service to get post's image
		var image string
		var err error
		getMediaRequest := dtos.GetMediaRequest{
			ID: uint64(postEntity.MediaID),
		}
		if image, err = s.mediaGrpcClient.GetMedia(getMediaRequest); err != nil {
			return nil, rest_error.NewInternalServerError("user grpc client error when getting media", err)
		}

		// GRPC CALL TO USER SERVICE FOR USERNAME
		var username string
		if username, err = s.userGrpcClient.GetUsername(dtos.GetUsernameRequest{Email: userEmail}); err != nil {
			return nil, rest_error.NewInternalServerError("user grpc client error when getting username", err)
		}

		// Check if logged user liked, disliked or added post to favorites
		liked := false
		disliked := false
		inFavorites := false
		if loggedInUserEmail != "" {
			if _, postErr = s.likesRepository.GetByUserAndPost(loggedInUserEmail, postEntity.ID); postErr == nil {
				liked = true
			}

			if _, postErr = s.dislikesRepository.GetByUserAndPost(loggedInUserEmail, postEntity.ID); postErr == nil {
				disliked = true
			}

			// GRPC CALL TO USER SERVICE TO CHECK IF POST IS IN USER'S FAVORITES
			checkFavoritesRequest := dtos.CheckFavoritesRequest{
				Email: loggedInUserEmail,
				PostID: postEntity.ID,
			}
			if inFavorites, err = s.userGrpcClient.CheckPostIsInFavorites(checkFavoritesRequest); err != nil {
				return nil, rest_error.NewInternalServerError("user grpc client error when checking favorites", err)
			}
		}

		// Calculate number of post's likes and dislikes
		var numberOfLikes int64
		if numberOfLikes, postErr = s.likesRepository.GetNumberOfLikes(postEntity.ID); postErr != nil {
			return nil, postErr
		}

		var numberOfDislikes int64
		if numberOfDislikes, postErr = s.dislikesRepository.GetNumberOfDislikes(postEntity.ID); postErr != nil {
			return nil, postErr
		}

		// Get posts's comments
		var commentsDTOs = make([]dtos.CommentDTO, 0)
		var comments []modelComment.Comment
		if comments, postErr = s.commentsRepository.GetComments(postEntity.ID); postErr != nil {
			return nil, postErr
		}
		for _, commentEntity := range comments {
			if username, err = s.userGrpcClient.GetUsername(dtos.GetUsernameRequest{Email: commentEntity.UserEmail}); err != nil {
				return nil, rest_error.NewInternalServerError("user grpc client error when getting username", err)
			}
			commentsDTOs = append(commentsDTOs, dtos.CommentDTO{
				Text:     s.ProcessTags(commentEntity.Text),
				Date:     time.Unix(commentEntity.Date, 0).Format(layout),
				Username: username,
			})
		}

		// CREATE POST DTO
		postsDTOs = append(postsDTOs, dtos.PostDTO{
			ID:          postEntity.ID,
			Description: description,
			Date:        date,
			Image:       image,
			Username:    username,
			Liked:       liked,
			Disliked:    disliked,
			InFavorites: inFavorites,
			Likes:       uint(numberOfLikes),
			Dislikes:    uint(numberOfDislikes),
			Comments:    commentsDTOs,
		})
	}

	return postsDTOs, nil
}

func (s *postsService) ProcessTags(text string) string {
	r := regexp.MustCompile(`@[A-Za-z0-9_.]+`)
	matches := r.FindAllString(text, -1)

	var link string
	var taggable bool
	checkTaggableRequest := dtos.CheckTaggableRequest{
		Username: "",
	}

	for _, tag:= range matches {
		checkTaggableRequest.Username = tag[1:]
		if taggable, _ = s.userGrpcClient.CheckIfUserIsTaggable(checkTaggableRequest); !taggable {
			continue
		}

		link = "<a href='/users/" + tag[1:] + "' >" + tag + "</a>"
		text = strings.ReplaceAll(text, tag, link)
	}

	return text
}