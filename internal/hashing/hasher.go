package hashing

type Hasher interface {
	CreateHash(v string) (string, error)
	ValidateHash(v, hash string) error
}
