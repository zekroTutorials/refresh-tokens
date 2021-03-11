package main

import (
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/zekroTutorials/refresh-tokens/internal/accesstokens"
)

type config struct {
	PublicKeyFile string `envconfig:"PUBLICKEYFILE" default:"./cert.pem"`

	BindAddress string `envconfig:"WS_BINDADDRESS" default:":8080"`
	Prefix      string `envconfig:"WS_PREFIX" default:""`
}

var (
	userSnowflakeNode          *snowflake.Node
	refreshTokensSnowflakeNode *snowflake.Node

	atvalidator accesstokens.Validator
)

func main() {
	var err error

	userSnowflakeNode, err = snowflake.NewNode(10)
	must(err)

	refreshTokensSnowflakeNode, err = snowflake.NewNode(20)
	must(err)

	cfg := new(config)
	envconfig.MustProcess("APP", cfg)

	atvalidator, err = accesstokens.NewJWTManager("", cfg.PublicKeyFile)
	must(err)

	r := gin.Default()

	r.Use(validateAccessToken)

	r.GET(cfg.Prefix+"/me", getMe)

	r.Run(cfg.BindAddress)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
