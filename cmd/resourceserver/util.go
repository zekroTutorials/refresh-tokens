package main

import "github.com/gin-gonic/gin"

type errorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func jsonError(ctx *gin.Context, code int, err error) {
	ctx.JSON(code, &errorResponse{code, err.Error()})
}
