package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSONResult(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"result": data,
	})
}

func JSONError(ctx *gin.Context, status int, msg string) {
	ctx.JSON(status, gin.H{
		"error": msg,
	})
}
