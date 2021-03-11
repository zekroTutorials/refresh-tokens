package main

import "github.com/gin-gonic/gin"

func getMe(ctx *gin.Context) {
	ident, _ := ctx.Get("ident")
	ctx.JSON(200, ident)
}
