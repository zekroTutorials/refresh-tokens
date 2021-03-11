package models

import "time"

// EntityModel holds generic database
// entity properties.
type EntityModel struct {
	ID      string    `json:"id"`
	Created time.Time `json:"created"`
}

// IsNil returns true when the underlying
// entity instance is nil.
func (e *EntityModel) IsNil() bool {
	return e == nil
}
