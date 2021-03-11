package accesstokens

type Generator interface {
	Generate(ident string) (token string, err error)
}

type Validator interface {
	Validate(token string) (ident string, err error)
}
