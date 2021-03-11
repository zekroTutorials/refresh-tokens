package hashing

// Hasher provides functionalities to create and
// compare secure hashes from passwords.
type Hasher interface {
	CreateHash(v string) (string, error)
	ValidateHash(v, hash string) error
}
