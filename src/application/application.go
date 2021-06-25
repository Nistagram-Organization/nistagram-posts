package application

import (
	"github.com/Nistagram-Organization/nistagram-posts/src/controllers/ping"
	controller "github.com/Nistagram-Organization/nistagram-posts/src/controllers/post"
	"github.com/Nistagram-Organization/nistagram-posts/src/datasources/mysql"
	postrepository "github.com/Nistagram-Organization/nistagram-posts/src/repositories/post"
	postservice "github.com/Nistagram-Organization/nistagram-posts/src/services/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/comment"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/dislike"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/like"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/user_tag"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	database := mysql.NewMySqlDatabaseClient()
	if err := database.Init(); err != nil {
		panic(err)
	}

	if err := database.Migrate(
		&post.Post{},
		&like.Like{},
		&dislike.Dislike{},
		&comment.Comment{},
		&user_tag.UserTag{},
	); err != nil {
		panic(err)
	}

	pingController := ping.NewPingController()
	postController := controller.PostController(
		postservice.NewPostService(
			postrepository.NewPostRepository(database),
		),
	)

	router.GET("/ping", pingController.Ping)
	router.GET("/posts", postController.GetAll)

	router.Run(":8085")
}
