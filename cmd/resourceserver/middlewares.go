package main

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zekroTutorials/refresh-tokens/internal/wsutil"
)

var errInvalidAccessToken = errors.New("invalid access token")

func validateAccessToken(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	if !strings.HasPrefix(authHeader, "accessToken ") {
		atAbort(ctx)
		return
	}

	accessToken := authHeader[12:]
	ident, err := atvalidator.Validate(accessToken)
	if err != nil {
		atAbort(ctx)
		return
	}

	ctx.Set("ident", ident)
}

func atAbort(ctx *gin.Context) {
	wsutil.JsonError(ctx, 401, errInvalidAccessToken)
	ctx.Abort()
}
