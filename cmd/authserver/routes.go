package main

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zekroTutorials/refresh-tokens/internal/models"
	"github.com/zekroTutorials/refresh-tokens/internal/wsutil"
	"github.com/zekroTutorials/refresh-tokens/pkg/random"
)

type loginRequestModel struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type accessTokenModel struct {
	Token    string    `json:"token"`
	Deadline time.Time `json:"deadline"`
}

const (
	sessionExpiration = 2 * time.Hour
	sessionCookieName = "refreshToken"
)

var errUnauthorized = errors.New("unauthorized")

///////////////////////////////////////////////////////////////
// POST /login
///////////////////////////////////////////////////////////////

func postLogin(ctx *gin.Context) {
	data := new(loginRequestModel)
	if err := json.NewDecoder(ctx.Request.Body).Decode(data); err != nil {
		wsutil.JsonError(ctx, 400, err)
		return
	}

	user, err := db.GetUser(data.UserName)
	if err != nil {
		wsutil.JsonError(ctx, 500, err)
		return
	}
	if user == nil {
		wsutil.JsonError(ctx, 401, errUnauthorized)
		return
	}

	if err = hasher.ValidateHash(data.Password, user.PasswordHash); err != nil {
		wsutil.JsonError(ctx, 401, errUnauthorized)
		return
	}

	tokenStr, err := random.Base64(64)
	if err != nil {
		wsutil.JsonError(ctx, 500, err)
		return
	}

	token := &models.RefreshToken{
		EntityModel: &models.EntityModel{
			ID:      refreshTokensSnowflakeNode.Generate().String(),
			Created: time.Now(),
		},
		UserID:   user.ID,
		Token:    tokenStr,
		Deadline: time.Now().Add(sessionExpiration),
	}

	if err = db.AddRefreshToken(token); err != nil {
		wsutil.JsonError(ctx, 500, err)
		return
	}

	ctx.SetCookie(sessionCookieName, tokenStr, int(sessionExpiration), "", "", false, true)

	ctx.JSON(200, user.Sanitize())
}

///////////////////////////////////////////////////////////////
// GET /accesstoken
///////////////////////////////////////////////////////////////

func getAccesstoken(ctx *gin.Context) {
	refreshToken, _ := ctx.Cookie(sessionCookieName)
	if refreshToken == "" {
		wsutil.JsonError(ctx, 401, errUnauthorized)
		return
	}

	rtModel, err := db.GetRefreshToken(refreshToken)
	if err != nil {
		wsutil.JsonError(ctx, 500, err)
		return
	}
	if rtModel.IsNil() || time.Now().After(rtModel.Deadline) {
		wsutil.JsonError(ctx, 401, errUnauthorized)
		return
	}

	accessToken, err := atgenerator.Generate(rtModel.UserID, accessTokenLifetime)
	if err != nil {
		wsutil.JsonError(ctx, 500, err)
		return
	}

	ctx.JSON(200, &accessTokenModel{accessToken, time.Now().Add(accessTokenLifetime)})
}
