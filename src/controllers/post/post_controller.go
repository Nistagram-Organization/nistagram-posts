package post

import (
	"github.com/Nistagram-Organization/agent-shared/src/utils/rest_error"
	"github.com/Nistagram-Organization/nistagram-posts/src/services/post"
	"github.com/Nistagram-Organization/nistagram-posts/src/dtos"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PostController interface {
	GetAll(*gin.Context)
	LikePost(*gin.Context)
}

type postsController struct {
	postsService post.PostService
}

func NewPostController(postsService post.PostService) PostController {
	return &postsController{
		postsService: postsService,
	}
}

func getId(idParam string) (uint, rest_error.RestErr) {
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return 0, rest_error.NewBadRequestError("Id should be a number")
	}
	return uint(id), nil
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

func (p *postsController) GetAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, p.postsService.GetAll())
}
