package post

import (
	"github.com/Nistagram-Organization/nistagram-posts/src/services/post"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PostController interface {
	GetAll(*gin.Context)
}

type postsController struct {
	postsService post.PostService
}

func NewPostController(postsService post.PostService) PostController {
	return &postsController{
		postsService: postsService,
	}
}

func (p *postsController) GetAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, p.postsService.GetAll())
}
