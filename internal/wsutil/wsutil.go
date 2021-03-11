package wsutil

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func JsonError(ctx *gin.Context, code int, err error) {
	ctx.JSON(code, &ErrorResponse{code, err.Error()})
}

func AddCorsHeader(publicAddr string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", publicAddr)
		ctx.Header("Access-Control-Allow-Headers", "authorization, content-type, set-cookie, cookie, server")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, POST, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if strings.ToLower(ctx.Request.Method) == "options" {
			ctx.Status(204)
			ctx.Abort()
		}
	}
}
