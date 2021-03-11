package models

import "time"

type RefreshToken struct {
	*EntityModel

	UserID   string
	Token    string
	Deadline time.Time
}
