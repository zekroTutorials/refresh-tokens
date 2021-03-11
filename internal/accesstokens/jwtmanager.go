package accesstokens

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const tokenLifetime = 30 * time.Minute

var jwtSigningMethod = jwt.SigningMethodRS256

type JWTManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewJWTManager(privateKeyFile, publicKeyFile string) (m *JWTManager, err error) {
	m = new(JWTManager)

	var kd []byte
	if privateKeyFile != "" {
		kd, err = ioutil.ReadFile(privateKeyFile)
		if err != nil {
			return
		}
		m.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(kd)
		if err != nil {
			return
		}
	}

	if publicKeyFile != "" {
		kd, err = ioutil.ReadFile(publicKeyFile)
		if err != nil {
			return
		}
		m.publicKey, err = jwt.ParseRSAPublicKeyFromPEM(kd)
		if err != nil {
			return
		}
	}

	return
}

func (m *JWTManager) Generate(ident string) (token string, err error) {
	if m.privateKey == nil {
		err = errors.New("not supported with this instance")
		return
	}

	now := time.Now()
	token, err = jwt.NewWithClaims(jwtSigningMethod, jwt.StandardClaims{
		Subject:   ident,
		ExpiresAt: now.Add(tokenLifetime).Unix(),
		IssuedAt:  now.Unix(),
	}).SignedString(m.privateKey)

	return
}

func (m *JWTManager) Validate(token string) (ident string, err error) {
	if m.publicKey == nil {
		err = errors.New("not supported with this instance")
		return
	}

	claims := new(jwt.StandardClaims)
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return m.publicKey, nil
	})
	if err != nil {
		return
	}

	if !jwtToken.Valid {
		err = errors.New("invalid claims")
	}

	ident = claims.Subject

	return
}
