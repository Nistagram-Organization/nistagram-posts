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
	"github.com/Nistagram-Organization/nistagram-posts/src/services/post_grpc_service"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/comment"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/dislike"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/like"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/user_tag"
	"github.com/Nistagram-Organization/nistagram-shared/src/proto"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/prometheus_handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	dockerKey = "docker"
)

var (
	router = gin.Default()
)

func StartApplication() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	router.Use(cors.New(corsConfig))

	var docker bool
	if os.Getenv(dockerKey) == "" {
		docker = false
	} else {
		docker = true
	}

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

	port := ":8085"
	l, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	m := cmux.New(l)

	grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	mediaGrpcClient := media_grpc_client.NewMediaGrpcClient(docker)
	userGrpcClient := user_grpc_client.NewUserGrpcClient(docker)
	commentRepo := commentRepository.NewCommentRepository(database)
	dislikeRepo := dislikerepository.NewDislikeRepository(database)
	likeRepo := likerepository.NewLikeRepository(database)
	postRepo := postrepository.NewPostRepository(database)
	postService := postservice.NewPostService(postRepo, likeRepo, dislikeRepo, commentRepo, mediaGrpcClient, userGrpcClient)
	postGrpcService := post_grpc_service.NewPostGrpcService(postService)

	postController := controller.NewPostController(postService)

	router.POST("/posts", postController.CreatePost)
	router.POST("/posts/like", postController.LikePost)
	router.DELETE("/posts/like", postController.UnlikePost)
	router.POST("/posts/dislike", postController.DislikePost)
	router.DELETE("/posts/dislike", postController.UndislikePost)
	router.POST("/posts/report/:id", postController.ReportInappropriateContent)
	router.POST("/posts/comment", postController.PostComment)
	router.GET("/posts", postController.GetUsersPosts)
	router.GET("/posts/inappropriate", postController.GetInappropriateContent)
	router.GET("/posts/feed", postController.GetPostsFeed)
	router.GET("/posts/search", postController.SearchTags)

	router.GET("/metrics", prometheus_handler.PrometheusGinHandler())

	grpcS := grpc.NewServer()
	proto.RegisterPostServiceServer(grpcS, postGrpcService)

	httpS := &http.Server{
		Handler: router,
	}

	go grpcS.Serve(grpcListener)
	go httpS.Serve(httpListener)

	log.Printf("Running http and grpc server on port %s", port)
	m.Serve()
}
