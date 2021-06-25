package ping

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type PingController interface {
	Ping(ctx *gin.Context)
}

type pingController struct {
}

func NewPingController() PingController {
	return &pingController{}
}

func (c *pingController) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "ping")
}
