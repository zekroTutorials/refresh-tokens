package database

import (
	"context"
	"time"

	"github.com/zekroTutorials/refresh-tokens/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDriver struct {
	client *mongo.Client

	users         *mongo.Collection
	refreshTokens *mongo.Collection
}

func NewMongoDriver(connectionString, database string) (m *MongoDriver, err error) {
	m = new(MongoDriver)

	if m.client, err = mongo.NewClient(options.Client().ApplyURI(connectionString)); err != nil {
		return
	}

	ctxConnect, cancelConnect := ctxTimeout(5 * time.Second)
	defer cancelConnect()

	if err = m.client.Connect(ctxConnect); err != nil {
		return
	}

	ctxPing, cancelPing := ctxTimeout(5 * time.Second)
	defer cancelPing()

	if err = m.client.Ping(ctxPing, readpref.Primary()); err != nil {
		return
	}

	db := m.client.Database(database)
	m.users = db.Collection("users")
	m.refreshTokens = db.Collection("refreshTokens")

	return
}

func (m *MongoDriver) AddUser(user *models.UserModel) (err error) {
	ctx, cancel := ctxTimeout(5 * time.Second)
	defer cancel()

	_, err = m.users.InsertOne(ctx, user)
	return
}

func (m *MongoDriver) GetUser(ident string) (user *models.UserModel, err error) {
	ctx, cancel := ctxTimeout(5 * time.Second)
	defer cancel()

	user = new(models.UserModel)
	err = m.users.FindOne(ctx, bson.M{
		"$or": bson.A{
			bson.M{"id": ident},
			bson.M{"username": ident},
		},
	}).Decode(user)

	err = wrapError(err)

	return
}

func (m *MongoDriver) AddRefreshToken(token *models.RefreshToken) (err error) {
	ctx, cancel := ctxTimeout(5 * time.Second)
	defer cancel()

	_, err = m.refreshTokens.InsertOne(ctx, token)
	return
}

func (m *MongoDriver) GetRefreshToken(token string) (res *models.RefreshToken, err error) {
	ctx, cancel := ctxTimeout(5 * time.Second)
	defer cancel()

	res = new(models.RefreshToken)
	err = m.refreshTokens.FindOne(ctx, bson.M{"token": token}).Decode(res)

	err = wrapError(err)

	return
}

func (m *MongoDriver) DeleteRefreshToken(id string) error {
	ctx, cancel := ctxTimeout(5 * time.Second)
	defer cancel()

	_, err := m.refreshTokens.DeleteOne(ctx, bson.M{"id": id})
	return wrapError(err)
}

func ctxTimeout(d time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	return ctx, cancel
}

func wrapError(err error) error {
	if err == nil || err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}
