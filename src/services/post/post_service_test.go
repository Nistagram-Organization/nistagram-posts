package post

import (
	"errors"
	"fmt"
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PostServiceUnitTestsSuite struct {
	suite.Suite
	postsRepositoryMock    *post.PostRepositoryMock
	likesRepositoryMock    *like.LikeRepositoryMock
	dislikesRepositoryMock *dislike.DislikeRepositoryMock
	commentsRepositoryMock *comment.CommentRepositoryMock
	mediaGrpcClientMock    *media_grpc_client.MediaGrpcClientMock
	userGrpcClient         *user_grpc_client.UserGrpcClient
	service                PostService
}

func TestPostServiceUnitTestsSuite(t *testing.T) {
	suite.Run(t, new(PostServiceUnitTestsSuite))
}

func (suite *PostServiceUnitTestsSuite) SetupSuite() {
	suite.postsRepositoryMock = new(post.PostRepositoryMock)
	suite.likesRepositoryMock = new(like.LikeRepositoryMock)
	suite.dislikesRepositoryMock = new(dislike.DislikeRepositoryMock)
	suite.commentsRepositoryMock = new(comment.CommentRepositoryMock)
	suite.mediaGrpcClientMock = new(media_grpc_client.MediaGrpcClientMock)
	suite.userGrpcClient = new(user_grpc_client.UserGrpcClient)
	suite.service = NewPostService(suite.postsRepositoryMock, suite.likesRepositoryMock, suite.dislikesRepositoryMock,
		suite.commentsRepositoryMock, suite.mediaGrpcClientMock, *suite.userGrpcClient)
}

func (suite *PostServiceUnitTestsSuite) TestNewPostService() {
	assert.NotNil(suite.T(), suite.service, "Service is nil")
}

func (suite *PostServiceUnitTestsSuite) TestPostService_LikePost_PostDoesNotExist() {
	likeRequestDTO := dtos.LikeDislikeRequestDTO{
		PostID:    1,
		UserEmail: "",
	}
	err := rest_error.NewNotFoundError(fmt.Sprintf("Error when trying to get post with id %d", likeRequestDTO.PostID))

	suite.postsRepositoryMock.On("Get", likeRequestDTO.PostID).Return(nil, err).Once()

	likeErr := suite.service.LikePost(&likeRequestDTO)

	assert.Equal(suite.T(), err, likeErr)
}

func (suite *PostServiceUnitTestsSuite) TestPostService_LikePost_PostAlreadyLiked() {
	likeRequestDTO := dtos.LikeDislikeRequestDTO{
		PostID:    1,
		UserEmail: "mail@mail.com",
	}
	err := rest_error.NewBadRequestError("Post already liked")

	suite.postsRepositoryMock.On("Get", likeRequestDTO.PostID).Return(&modelPost.Post{}, nil).Once()
	suite.likesRepositoryMock.On("GetByUserAndPost", likeRequestDTO.UserEmail, likeRequestDTO.PostID).Return(&modelLike.Like{}, nil).Once()

	likeErr := suite.service.LikePost(&likeRequestDTO)

	assert.Equal(suite.T(), err, likeErr)
}

func (suite *PostServiceUnitTestsSuite) TestPostService_LikePost_PostAlreadyDisliked() {
	likeRequestDTO := dtos.LikeDislikeRequestDTO{
		PostID:    1,
		UserEmail: "mail@mail.com",
	}
	err := rest_error.NewBadRequestError("Post already disliked")

	suite.postsRepositoryMock.On("Get", likeRequestDTO.PostID).Return(&modelPost.Post{}, nil).Once()
	suite.likesRepositoryMock.On("GetByUserAndPost", likeRequestDTO.UserEmail, likeRequestDTO.PostID).Return(&modelLike.Like{}, err).Once()
	suite.dislikesRepositoryMock.On("GetByUserAndPost", likeRequestDTO.UserEmail, likeRequestDTO.PostID).Return(&modelDislike.Dislike{}, nil).Once()

	likeErr := suite.service.LikePost(&likeRequestDTO)

	assert.Equal(suite.T(), err, likeErr)
}

func (suite *PostServiceUnitTestsSuite) TestPostService_LikePost() {
	likeRequestDTO := dtos.LikeDislikeRequestDTO{
		PostID:    1,
		UserEmail: "mail@mail.com",
	}
	likeEntity := modelLike.Like{
		UserEmail: likeRequestDTO.UserEmail,
		PostID:    likeRequestDTO.PostID,
	}
	err := rest_error.NewBadRequestError("Post already disliked")

	suite.postsRepositoryMock.On("Get", likeRequestDTO.PostID).Return(&modelPost.Post{}, nil).Once()
	suite.likesRepositoryMock.On("GetByUserAndPost", likeRequestDTO.UserEmail, likeRequestDTO.PostID).Return(&modelLike.Like{}, err).Once()
	suite.dislikesRepositoryMock.On("GetByUserAndPost", likeRequestDTO.UserEmail, likeRequestDTO.PostID).Return(&modelDislike.Dislike{}, err).Once()
	suite.likesRepositoryMock.On("Create", &likeEntity).Return(nil)

	likeErr := suite.service.LikePost(&likeRequestDTO)

	assert.Equal(suite.T(), nil, likeErr)
}

func (suite *PostServiceUnitTestsSuite) TestPostService_DislikePost_PostDoesNotExist() {
	dislikeRequestDTO := dtos.LikeDislikeRequestDTO{
		PostID:    1,
		UserEmail: "",
	}
	err := rest_error.NewNotFoundError(fmt.Sprintf("Error when trying to get post with id %d", dislikeRequestDTO.PostID))

	suite.postsRepositoryMock.On("Get", dislikeRequestDTO.PostID).Return(nil, err).Once()

	dislikeErr := suite.service.DislikePost(&dislikeRequestDTO)

	assert.Equal(suite.T(), err, dislikeErr)
}

func (suite *PostServiceUnitTestsSuite) TestPostService_DislikePost_PostAlreadyDisliked() {
	dislikeRequestDTO := dtos.LikeDislikeRequestDTO{
		PostID:    1,
		UserEmail: "mail@mail.com",
	}
	err := rest_error.NewBadRequestError("Post already disliked")

	suite.postsRepositoryMock.On("Get", dislikeRequestDTO.PostID).Return(&modelPost.Post{}, nil).Once()
	suite.dislikesRepositoryMock.On("GetByUserAndPost", dislikeRequestDTO.UserEmail, dislikeRequestDTO.PostID).Return(&modelDislike.Dislike{}, nil).Once()

	dislikeErr := suite.service.DislikePost(&dislikeRequestDTO)

	assert.Equal(suite.T(), err, dislikeErr)
}

func (suite *PostServiceUnitTestsSuite) TestPostService_DislikePost_PostAlreadyLiked() {
	dislikeRequestDTO := dtos.LikeDislikeRequestDTO{
		PostID:    1,
		UserEmail: "mail@mail.com",
	}
	err := rest_error.NewBadRequestError("Post already liked")

	suite.postsRepositoryMock.On("Get", dislikeRequestDTO.PostID).Return(&modelPost.Post{}, nil).Once()
	suite.dislikesRepositoryMock.On("GetByUserAndPost", dislikeRequestDTO.UserEmail, dislikeRequestDTO.PostID).Return(&modelDislike.Dislike{}, err).Once()
	suite.likesRepositoryMock.On("GetByUserAndPost", dislikeRequestDTO.UserEmail, dislikeRequestDTO.PostID).Return(&modelLike.Like{}, nil).Once()

	dislikeErr := suite.service.DislikePost(&dislikeRequestDTO)

	assert.Equal(suite.T(), err, dislikeErr)
}

func (suite *PostServiceUnitTestsSuite) TestPostService_DislikePost() {
	dislikeRequestDTO := dtos.LikeDislikeRequestDTO{
		PostID:    1,
		UserEmail: "mail@mail.com",
	}
	dislikeEntity := modelDislike.Dislike{
		UserEmail: dislikeRequestDTO.UserEmail,
		PostID:    dislikeRequestDTO.PostID,
	}
	err := rest_error.NewBadRequestError("Post already disliked")

	suite.postsRepositoryMock.On("Get", dislikeRequestDTO.PostID).Return(&modelPost.Post{}, nil).Once()
	suite.dislikesRepositoryMock.On("GetByUserAndPost", dislikeRequestDTO.UserEmail, dislikeRequestDTO.PostID).Return(&modelDislike.Dislike{}, err).Once()
	suite.likesRepositoryMock.On("GetByUserAndPost", dislikeRequestDTO.UserEmail, dislikeRequestDTO.PostID).Return(&modelLike.Like{}, err).Once()
	suite.dislikesRepositoryMock.On("Create", &dislikeEntity).Return(nil)

	dislikeErr := suite.service.DislikePost(&dislikeRequestDTO)

	assert.Equal(suite.T(), nil, dislikeErr)
}

func (suite *PostServiceUnitTestsSuite) TestPostService_PostComment_PostDoesNotExist() {
	commentEntity := modelComment.Comment{
		PostID: 1,
	}
	err := rest_error.NewNotFoundError(fmt.Sprintf("Error when trying to get post with id %d", commentEntity.PostID))

	suite.postsRepositoryMock.On("Get", commentEntity.PostID).Return(nil, err).Once()

	commErr := suite.service.PostComment(&commentEntity)

	assert.Equal(suite.T(), err, commErr)
}

func (suite *PostServiceUnitTestsSuite) TestPostService_PostComment() {
	commentEntity := modelComment.Comment{
		PostID: 1,
	}

	suite.postsRepositoryMock.On("Get", commentEntity.PostID).Return(&modelPost.Post{}, nil).Once()
	suite.commentsRepositoryMock.On("Create", &commentEntity).Return(nil).Once()

	commErr := suite.service.PostComment(&commentEntity)

	assert.Equal(suite.T(), nil, commErr)
}

func (suite *PostServiceUnitTestsSuite) TestPostService_CreatePost_GRPCError() {
	postDTO := dtos.CreatePostDTO{
		Description: "Opis",
		Image:       "Image",
		UserEmail:   "mail@mail.com",
	}
	saveMediaRequest := dtos.SaveMediaRequest{
		Image: "Image",
	}
	err := rest_error.NewInternalServerError("user grpc client error when saving media", errors.New(""))

	suite.mediaGrpcClientMock.On("SaveMedia", saveMediaRequest).Return(new(uint), err).Once()

	createErr := suite.service.CreatePost(&postDTO)

	assert.Equal(suite.T(), err, createErr)
}

func (suite *PostServiceUnitTestsSuite) TestPostService_CreatePost() {
	postDTO := dtos.CreatePostDTO{
		Description: "Opis",
		Image:       "Image",
		UserEmail:   "mail@mail.com",
	}
	saveMediaRequest := dtos.SaveMediaRequest{
		Image: "Image",
	}
	postEntity := modelPost.Post{
		Description:           postDTO.Description,
		UserEmail:             postDTO.UserEmail,
		MarkedAsInappropriate: false,
		Date:                  time_utils.Now(),
		MediaID:               0,
	}
	err := rest_error.NewInternalServerError("user grpc client error when saving media", nil)

	suite.mediaGrpcClientMock.On("SaveMedia", saveMediaRequest).Return(new(uint), err).Once()
	suite.postsRepositoryMock.On("Create", &postEntity).Return(nil).Once()

	createErr := suite.service.CreatePost(&postDTO)

	assert.Equal(suite.T(), nil, createErr)
}
