package post

import (
	"fmt"
	"github.com/Nistagram-Organization/nistagram-posts/src/clients/media_grpc_client"
	"github.com/Nistagram-Organization/nistagram-posts/src/clients/user_grpc_client"
	"github.com/Nistagram-Organization/nistagram-posts/src/datasources/mysql"
	"github.com/Nistagram-Organization/nistagram-posts/src/dtos"
	commentRepository "github.com/Nistagram-Organization/nistagram-posts/src/repositories/comment"
	dislikerepository "github.com/Nistagram-Organization/nistagram-posts/src/repositories/dislike"
	likerepository "github.com/Nistagram-Organization/nistagram-posts/src/repositories/like"
	postrepository "github.com/Nistagram-Organization/nistagram-posts/src/repositories/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/comment"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/dislike"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/like"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/user_tag"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"os"
	"testing"
)

type PostServiceIntegrationTestsSuite struct {
	suite.Suite
	service PostService
	db      *gorm.DB
	posts   []post.Post
}

func (suite *PostServiceIntegrationTestsSuite) SetupSuite() {
	var docker bool
	if os.Getenv("docker") == "" {
		docker = false
	} else {
		docker = true
	}

	database := mysql.NewMySqlDatabaseClient()
	if err := database.Init(); err != nil {
		suite.Fail("Failed to initialize database")
	}

	if err := database.Migrate(
		&like.Like{},
		&dislike.Dislike{},
		&comment.Comment{},
		&user_tag.UserTag{},
		&post.Post{},
	); err != nil {
		panic(err)
	}

	suite.db = database.GetClient()

	mediaGrpcClient := media_grpc_client.NewMediaGrpcClient(false)
	userGrpcClient := user_grpc_client.NewUserGrpcClient(docker)
	commentRepo := commentRepository.NewCommentRepository(database)
	dislikeRepo := dislikerepository.NewDislikeRepository(database)
	likeRepo := likerepository.NewLikeRepository(database)
	postRepo := postrepository.NewPostRepository(database)
	suite.service = NewPostService(postRepo, likeRepo, dislikeRepo, commentRepo, mediaGrpcClient, userGrpcClient)
}

func (suite *PostServiceIntegrationTestsSuite) SetupTest() {
	suite.posts = []post.Post{
		{
			ID:                    1,
			Description:           "Opis",
			Date:                  123456,
			MarkedAsInappropriate: false,
			UserEmail:             "mail@mail.com",
			MediaID:               1,
		},
		{
			ID:                    2,
			Description:           "Opis",
			Date:                  1234,
			MarkedAsInappropriate: false,
			UserEmail:             "mail@mail.com",
			MediaID:               2,
		},
		{
			ID:                    3,
			Description:           "Opis",
			Date:                  12345,
			MarkedAsInappropriate: false,
			UserEmail:             "mail@mail.com",
			MediaID:               3,
		},
		{
			ID:                    4,
			Description:           "Opis",
			Date:                  12345,
			MarkedAsInappropriate: false,
			UserEmail:             "mail@mail.com",
			MediaID:               4,
		},
	}
	likeEntity := like.Like{
		ID:        1,
		UserEmail: "mail@mail.com",
		PostID:    3,
	}
	dislikeEntity := dislike.Dislike{
		ID:        1,
		UserEmail: "mail@mail.com",
		PostID:    4,
	}

	tx := suite.db.Begin()
	tx.Create(&suite.posts[0])
	tx.Create(&suite.posts[1])
	tx.Create(&suite.posts[2])
	tx.Create(&suite.posts[3])
	tx.Create(likeEntity)
	tx.Create(dislikeEntity)
	tx.Commit()
}

func (suite *PostServiceIntegrationTestsSuite) TearDownTest() {
	tx := suite.db.Begin()
	session := &gorm.Session{AllowGlobalUpdate: true}
	tx.Session(session).Delete(&post.Post{})
	tx.Commit()
}

func TestPostServiceIntegrationTestsSuite(t *testing.T) {
	suite.Run(t, new(PostServiceIntegrationTestsSuite))
}

func (suite *PostServiceIntegrationTestsSuite) TestIntegrationPostService_LikePost_PostDoesNotExist() {
	likeRequestDTO := dtos.LikeDislikeRequestDTO{
		PostID:    10000,
		UserEmail: "mail@mail.com",
	}
	err := rest_error.NewNotFoundError(fmt.Sprintf("Error when trying to get post with id %d", likeRequestDTO.PostID))

	likeErr := suite.service.LikePost(&likeRequestDTO)

	assert.Equal(suite.T(), err, likeErr)
}

func (suite *PostServiceIntegrationTestsSuite) TestIntegrationPostService_LikePost() {
	likeRequestDTO := dtos.LikeDislikeRequestDTO{
		PostID:    2,
		UserEmail: "mail@mail.com",
	}

	likeErr := suite.service.LikePost(&likeRequestDTO)

	assert.Equal(suite.T(), nil, likeErr)
}

func (suite *PostServiceIntegrationTestsSuite) TestIntegrationPostService_UnlikePost_PostDoesNotExist() {
	id := uint(10000)
	err := rest_error.NewNotFoundError(fmt.Sprintf("Error when trying to get post with id %d", id))

	likeErr := suite.service.UnlikePost("mail@mail.com", id)

	assert.Equal(suite.T(), err, likeErr)
}

func (suite *PostServiceIntegrationTestsSuite) TestIntegrationPostService_UnlikePost() {
	likeErr := suite.service.UnlikePost("mail@mail.com", 3)

	assert.Equal(suite.T(), nil, likeErr)
}

func (suite *PostServiceIntegrationTestsSuite) TestIntegrationPostService_UndislikePost_PostDoesNotExist() {
	id := uint(10000)
	err := rest_error.NewNotFoundError(fmt.Sprintf("Error when trying to get post with id %d", id))

	likeErr := suite.service.UndislikePost("mail@mail.com", id)

	assert.Equal(suite.T(), err, likeErr)
}

func (suite *PostServiceIntegrationTestsSuite) TestIntegrationPostService_UndislikePost() {
	likeErr := suite.service.UndislikePost("mail@mail.com", 4)

	assert.Equal(suite.T(), nil, likeErr)
}

func (suite *PostServiceIntegrationTestsSuite) TestIntegrationPostService_DislikePost_PostDoesNotExist() {
	dislikeRequestDTO := dtos.LikeDislikeRequestDTO{
		PostID:    10000,
		UserEmail: "mail@mail.com",
	}
	err := rest_error.NewNotFoundError(fmt.Sprintf("Error when trying to get post with id %d", dislikeRequestDTO.PostID))

	dislikeErr := suite.service.DislikePost(&dislikeRequestDTO)

	assert.Equal(suite.T(), err, dislikeErr)
}

func (suite *PostServiceIntegrationTestsSuite) TestIntegrationPostService_DislikePost() {
	dislikeRequestDTO := dtos.LikeDislikeRequestDTO{
		PostID:    1,
		UserEmail: "mail@mail.com",
	}

	dislikeErr := suite.service.DislikePost(&dislikeRequestDTO)

	assert.Equal(suite.T(), nil, dislikeErr)
}

func (suite *PostServiceIntegrationTestsSuite) TestIntegrationPostService_PostComment_PostDoesNotExist() {
	commentEntity := comment.Comment{
		PostID: 10000,
	}
	err := rest_error.NewNotFoundError(fmt.Sprintf("Error when trying to get post with id %d", commentEntity.PostID))

	commErr := suite.service.PostComment(&commentEntity)

	assert.Equal(suite.T(), err, commErr)
}

func (suite *PostServiceIntegrationTestsSuite) TestIntegrationPostService_PostComment() {
	commentEntity := comment.Comment{
		PostID: 1,
	}

	commErr := suite.service.PostComment(&commentEntity)

	assert.Equal(suite.T(), nil, commErr)
}

func (suite *PostServiceIntegrationTestsSuite) TestIntegrationPostService_ReportInappropriatePost_PostDoesNotExist() {
	id := uint(10000)
	err := rest_error.NewNotFoundError(fmt.Sprintf("Error when trying to get post with id %d", id))

	reportErr := suite.service.ReportInappropriateContent(id)

	assert.Equal(suite.T(), err, reportErr)
}

func (suite *PostServiceIntegrationTestsSuite) TestIntegrationPostService_ReportInappropriatePost() {
	reportErr := suite.service.ReportInappropriateContent(1)

	assert.Equal(suite.T(), nil, reportErr)
}

func (suite *PostServiceIntegrationTestsSuite) TestIntegrationPostService_DecideOnPost_PostDoesNotExist() {
	id := uint(10000)
	err := rest_error.NewNotFoundError(fmt.Sprintf("Error when trying to get post with id %d", id))

	reportErr := suite.service.DecideOnContent(id, true)

	assert.Equal(suite.T(), err, reportErr)
}

func (suite *PostServiceIntegrationTestsSuite) TestIntegrationPostService_DecideOnPost() {
	decideErr := suite.service.DecideOnContent(1, false)

	assert.Equal(suite.T(), nil, decideErr)
}
