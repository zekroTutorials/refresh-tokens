package models

import "time"

type EntityModel struct {
	ID      string    `json:"id"`
	Created time.Time `json:"created"`
}

func (e *EntityModel) IsNil() bool {
	return e == nil
}
