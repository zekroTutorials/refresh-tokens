package wsutil

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// ErrorResponse wraps an error
// response model.
type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

// JsonError wraps the passed error into an
// ErrorModel and responds it to the given
// request ctx with the passed code.
func JsonError(ctx *gin.Context, code int, err error) {
	ctx.JSON(code, &ErrorResponse{code, err.Error()})
}

// AddCorsHeader returns a new middleware function
// which adds CORS headers to OPTIONS requests
// and responds with status code 204.
func AddCorsHeader(publicAddr string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.ToLower(ctx.Request.Method) == "options" {
			ctx.Header("Access-Control-Allow-Origin", publicAddr)
			ctx.Header("Access-Control-Allow-Headers", "authorization, content-type, set-cookie, cookie, server")
			ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, POST, DELETE, OPTIONS")
			ctx.Header("Access-Control-Allow-Credentials", "true")
			ctx.Status(204)
			ctx.Abort()
		}
	}
}
