package hashing

import (
	"errors"
	"runtime"

	"github.com/alexedwards/argon2id"
)

// Argon2IDHasher implements Hasher using the
// Argon2ID algorithm.
type Argon2IDHasher struct {
	params *argon2id.Params
}

// NewArgon2IDHasher initializes a new instance of
// Argon2IDHasher with sime default configuration.
func NewArgon2IDHasher() (a *Argon2IDHasher) {
	a = new(Argon2IDHasher)

	cpus := runtime.NumCPU()

	a.params = &argon2id.Params{
		Memory:      128 * 1024,
		Iterations:  4,
		Parallelism: uint8(cpus),
		SaltLength:  16,
		KeyLength:   32,
	}

	return
}

func (a *Argon2IDHasher) CreateHash(v string) (string, error) {
	return argon2id.CreateHash(v, a.params)
}

func (a *Argon2IDHasher) ValidateHash(v, hash string) (err error) {
	ok, err := argon2id.ComparePasswordAndHash(v, hash)
	if err != nil {
		return
	}
	if !ok {
		err = errors.New("invalid hash")
	}
	return
}
