package models

// UserModel defines a user database
// model.
type UserModel struct {
	*EntityModel

	UserName     string `json:"username"`
	PasswordHash string `json:"passwordhash,omitempty"`
}

// Sanitize removes crucial information
// from the user model like the password
// hash.
func (u *UserModel) Sanitize() *UserModel {
	nu := *u
	nu.PasswordHash = ""
	return &nu
}
