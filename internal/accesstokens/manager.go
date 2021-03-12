package accesstokens

import "time"

// Generator provides functionalities to generate
// access tokens linking to the given ident with
// the given expiration duration.
type Generator interface {
	Generate(ident string, expire time.Duration) (token string, err error)
}

// Validator provides functionalities to validate
// a given token and recovering the passed ident
// on generation.
//
// When the given token is invalid, an error is
// returned.
type Validator interface {
	Validate(token string) (ident string, err error)
}
