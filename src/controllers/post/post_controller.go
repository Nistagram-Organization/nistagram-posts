package post

import (
	"github.com/Nistagram-Organization/nistagram-posts/src/dtos"
	"github.com/Nistagram-Organization/nistagram-posts/src/services/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PostController interface {
	GetAll(*gin.Context)
	LikePost(*gin.Context)
	UnlikePost(*gin.Context)
}

type postsController struct {
	postsService post.PostService
}

func NewPostController(postsService post.PostService) PostController {
	return &postsController{
		postsService: postsService,
	}
}

func (p *postsController) LikePost(ctx *gin.Context) {
	var likeRequest dtos.LikeRequestDTO
	if err := ctx.ShouldBindJSON(&likeRequest); err != nil {
		restErr := rest_error.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	ctx.JSON(http.StatusOK, p.postsService.LikePost(&likeRequest))
}

func (p *postsController) UnlikePost(ctx *gin.Context) {
	var likeRequest dtos.LikeRequestDTO
	if err := ctx.ShouldBindJSON(&likeRequest); err != nil {
		restErr := rest_error.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	ctx.JSON(http.StatusOK, p.postsService.UnlikePost(&likeRequest))
}

func (p *postsController) GetAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, p.postsService.GetAll())
}
