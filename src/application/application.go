package application

import (
	"github.com/Nistagram-Organization/nistagram-posts/src/clients/media_grpc_client"
	"github.com/Nistagram-Organization/nistagram-posts/src/clients/user_grpc_client"
	controller "github.com/Nistagram-Organization/nistagram-posts/src/controllers/post"
	"github.com/Nistagram-Organization/nistagram-posts/src/datasources/mysql"
	commentRepository "github.com/Nistagram-Organization/nistagram-posts/src/repositories/comment"
	dislikerepository "github.com/Nistagram-Organization/nistagram-posts/src/repositories/dislike"
	likerepository "github.com/Nistagram-Organization/nistagram-posts/src/repositories/like"
	postrepository "github.com/Nistagram-Organization/nistagram-posts/src/repositories/post"
	postservice "github.com/Nistagram-Organization/nistagram-posts/src/services/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/comment"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/dislike"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/like"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/user_tag"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	router.Use(cors.New(corsConfig))

	database := mysql.NewMySqlDatabaseClient()
	if err := database.Init(); err != nil {
		panic(err)
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

	mediaGrpcClient := media_grpc_client.NewMediaGrpcClient()
	userGrpcClient := user_grpc_client.NewUserGrpcClient()
	commentRepo := commentRepository.NewCommentRepository(database)
	dislikeRepo := dislikerepository.NewDislikeRepository(database)
	likeRepo := likerepository.NewLikeRepository(database)
	postRepo := postrepository.NewPostRepository(database)
	postService := postservice.NewPostService(postRepo, likeRepo, dislikeRepo, commentRepo, mediaGrpcClient, userGrpcClient)

	postController := controller.NewPostController(postService)

	router.POST("/posts", postController.CreatePost)
	router.POST("/posts/like", postController.LikePost)
	router.DELETE("/posts/like", postController.UnlikePost)
	router.POST("/posts/dislike", postController.DislikePost)
	router.DELETE("/posts/dislike", postController.UndislikePost)
	router.POST("/posts/report/:id", postController.ReportInappropriateContent)
	router.POST("/posts/comment", postController.PostComment)
	router.GET("/posts", postController.GetUsersPosts)

	router.Run(":8085")
}
