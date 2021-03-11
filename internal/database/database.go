package database

import "github.com/zekroTutorials/refresh-tokens/internal/models"

// Database provides functionalities to bind
// to a database provider.
type Database interface {
	AddUser(user *models.UserModel) error
	GetUser(ident string) (*models.UserModel, error)

	AddRefreshToken(token *models.RefreshToken) error
	GetRefreshToken(token string) (*models.RefreshToken, error)
	DeleteRefreshToken(id string) error
}
