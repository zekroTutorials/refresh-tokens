package models

type UserModel struct {
	*EntityModel

	UserName     string `json:"username"`
	PasswordHash string `json:"passwordhash,omitempty"`
}

func (u *UserModel) Sanitize() *UserModel {
	nu := *u
	nu.PasswordHash = ""
	return &nu
}
