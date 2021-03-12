package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/zekroTutorials/refresh-tokens/internal/accesstokens"
	"github.com/zekroTutorials/refresh-tokens/internal/database"
	"github.com/zekroTutorials/refresh-tokens/internal/hashing"
	"github.com/zekroTutorials/refresh-tokens/internal/models"
	"github.com/zekroTutorials/refresh-tokens/internal/wsutil"
)

type config struct {
	MongoConnectionString string `envconfig:"MONGO_CONNECTIONSTRING"`
	MongoDatabase         string `envconfig:"MONGO_DATABASE" default:"tokenapp"`

	FirstUserName     string `envconfig:"ROOT_USERNAME"`
	FirstUserPassword string `envconfig:"ROOT_PASSWORD"`

	PrivateKeyFile string `envconfig:"PRIVATEKEYFILE" default:"./key.pem"`

	BindAddress string `envconfig:"WS_BINDADDRESS" default:":8080"`
	Prefix      string `envconfig:"WS_PREFIX" default:""`
	PublicAddr  string `envconfig:"WS_PUBLICADDR"`
}

// The lifetime of an access token.
const accessTokenLifetime = 30 * time.Minute

var (
	userSnowflakeNode          *snowflake.Node
	refreshTokensSnowflakeNode *snowflake.Node

	db          database.Database
	hasher      hashing.Hasher
	atgenerator accesstokens.Generator
)

func main() {
	var err error

	userSnowflakeNode, err = snowflake.NewNode(10)
	must(err)

	refreshTokensSnowflakeNode, err = snowflake.NewNode(20)
	must(err)

	cfg := new(config)
	envconfig.MustProcess("APP", cfg)

	atgenerator, err = accesstokens.NewJWTManager(cfg.PrivateKeyFile, "")
	must(err)

	db, err = database.NewMongoDriver(cfg.MongoConnectionString, cfg.MongoDatabase)
	must(err)

	hasher = hashing.NewArgon2IDHasher()

	mustInitFirstUser(cfg.FirstUserName, cfg.FirstUserPassword)

	r := gin.Default()

	if cfg.PublicAddr != "" {
		r.Use(wsutil.AddCorsHeader(cfg.PublicAddr))
	}

	r.POST(cfg.Prefix+"/login", postLogin)
	r.GET(cfg.Prefix+"/accesstoken", getAccesstoken)

	r.Run(cfg.BindAddress)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func mustInitFirstUser(username, password string) {
	u, err := db.GetUser(username)
	must(err)

	if !u.IsNil() {
		fmt.Println(u)
		return
	}

	pwHash, err := hasher.CreateHash(password)
	must(err)

	must(db.AddUser(&models.UserModel{
		EntityModel: &models.EntityModel{
			ID:      userSnowflakeNode.Generate().String(),
			Created: time.Now(),
		},
		UserName:     username,
		PasswordHash: pwHash,
	}))
}
