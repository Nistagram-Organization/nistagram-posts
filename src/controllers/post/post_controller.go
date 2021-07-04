package post

import (
	"github.com/Nistagram-Organization/nistagram-posts/src/dtos"
	"github.com/Nistagram-Organization/nistagram-posts/src/services/post"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/comment"
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
	ReportInappropriateContent(*gin.Context)
	PostComment(*gin.Context)
	CreatePost(*gin.Context)
	GetUsersPost(ctx *gin.Context)
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

	likeErr := p.postsService.LikePost(&likeRequest)
	if likeErr != nil {
		ctx.JSON(likeErr.Status(), likeErr)
		return
	}

	ctx.JSON(http.StatusOK, likeErr)
}

func (p *postsController) UnlikePost(ctx *gin.Context) {
	postId, idErr := getId(ctx.Query("post_id"))
	if idErr != nil {
		ctx.JSON(idErr.Status(), idErr)
		return
	}

	unlikeErr := p.postsService.UnlikePost(ctx.Query("user_mail"), postId)
	if unlikeErr != nil {
		ctx.JSON(unlikeErr.Status(), unlikeErr)
		return
	}

	ctx.JSON(http.StatusOK, unlikeErr)
}

func (p *postsController) DislikePost(ctx *gin.Context) {
	var dislikeRequest dtos.LikeDislikeRequestDTO
	if err := ctx.ShouldBindJSON(&dislikeRequest); err != nil {
		restErr := rest_error.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	dislikeErr := p.postsService.DislikePost(&dislikeRequest)
	if dislikeErr != nil {
		ctx.JSON(dislikeErr.Status(), dislikeErr)
		return
	}

	ctx.JSON(http.StatusOK, dislikeErr)
}

func (p *postsController) UndislikePost(ctx *gin.Context) {
	postId, idErr := getId(ctx.Query("post_id"))
	if idErr != nil {
		ctx.JSON(idErr.Status(), idErr)
		return
	}

	undislikeErr := p.postsService.UndislikePost(ctx.Query("user_mail"), postId)
	if undislikeErr != nil {
		ctx.JSON(undislikeErr.Status(), undislikeErr)
		return
	}

	ctx.JSON(http.StatusOK, undislikeErr)
}

func (p *postsController) ReportInappropriateContent(ctx *gin.Context) {
	postId, idErr := getId(ctx.Param("id"))
	if idErr != nil {
		ctx.JSON(idErr.Status(), idErr)
		return
	}

	reportErr := p.postsService.ReportInappropriateContent(postId)
	if reportErr != nil {
		ctx.JSON(reportErr.Status(), reportErr)
		return
	}

	ctx.JSON(http.StatusOK, reportErr)
}

func (p *postsController) PostComment(ctx *gin.Context) {
	var commentEntity comment.Comment
	if err := ctx.ShouldBindJSON(&commentEntity); err != nil {
		restErr := rest_error.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	commentErr := p.postsService.PostComment(&commentEntity)
	if commentErr != nil {
		ctx.JSON(commentErr.Status(), commentErr)
		return
	}

	ctx.JSON(http.StatusOK, commentErr)
}

func (p *postsController) CreatePost(ctx *gin.Context) {
	var createPostDTO dtos.CreatePostDTO
	if err := ctx.ShouldBindJSON(&createPostDTO); err != nil {
		restErr := rest_error.NewBadRequestError("invalid json body")
		ctx.JSON(restErr.Status(), restErr)
		return
	}

	createErr := p.postsService.CreatePost(&createPostDTO)
	if createErr != nil {
		ctx.JSON(createErr.Status(), createErr)
		return
	}

	ctx.JSON(http.StatusOK, createErr)
}

func (p *postsController) GetAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, p.postsService.GetAll())
}

func (p *postsController) GetUsersPosts(ctx *gin.Context) {
	getErr := p.postsService.GetUsersPosts(ctx.Param("user"), ctx.Param("logged_in_user"))
	if getErr != nil {
		ctx.JSON(getErr.Status(), getErr)
		return
	}

	ctx.JSON(http.StatusOK, getErr)
}