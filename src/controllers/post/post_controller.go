package post

import (
	"github.com/Nistagram-Organization/nistagram-posts/src/dtos"
	"github.com/Nistagram-Organization/nistagram-posts/src/services/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PostController interface {
	GetAll(*gin.Context)
	LikePost(*gin.Context)
	UnlikePost(*gin.Context)
	DislikePost(ctx *gin.Context)
	UndislikePost(ctx *gin.Context)
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
	var likeRequest dtos.LikeDislikeRequestDTO
	if err := ctx.ShouldBindJSON(&likeRequest); err != nil {
		restErr := rest_error.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	ctx.JSON(http.StatusOK, p.postsService.LikePost(&likeRequest))
}

func (p *postsController) UnlikePost(ctx *gin.Context) {
	userId, idErr := getId(ctx.Query("user_id"))
	if idErr != nil {
		ctx.JSON(idErr.Status(), idErr)
		return
	}

	postId, idErr := getId(ctx.Query("post_id"))
	if idErr != nil {
		ctx.JSON(idErr.Status(), idErr)
		return
	}

	ctx.JSON(http.StatusOK, p.postsService.UnlikePost(userId, postId))
}

func (p *postsController) DislikePost(ctx *gin.Context) {
	var dislikeRequest dtos.LikeDislikeRequestDTO
	if err := ctx.ShouldBindJSON(&dislikeRequest); err != nil {
		restErr := rest_error.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	ctx.JSON(http.StatusOK, p.postsService.DislikePost(&dislikeRequest))
}

func (p *postsController) UndislikePost(ctx *gin.Context) {
	userId, idErr := getId(ctx.Query("user_id"))
	if idErr != nil {
		ctx.JSON(idErr.Status(), idErr)
		return
	}

	postId, idErr := getId(ctx.Query("post_id"))
	if idErr != nil {
		ctx.JSON(idErr.Status(), idErr)
		return
	}

	ctx.JSON(http.StatusOK, p.postsService.UndislikePost(userId, postId))
}

func (p *postsController) GetAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, p.postsService.GetAll())
}
