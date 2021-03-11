package models

import "time"

// RefreshToken defines the database structure
// of a refresh token.
type RefreshToken struct {
	*EntityModel

	UserID   string
	Token    string
	Deadline time.Time
}
