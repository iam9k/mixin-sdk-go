package mixin

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type SignClaims struct {
	Uid        string
	Sid        string
	PrivateKey string
	Method     string
	Uri        string
	Body       string
	Scope      string
	Expire     int64
}

func SignAuthenticationToken(claims SignClaims) (string, error) {
	sum := sha256.Sum256([]byte(claims.Method + claims.Uri + claims.Body))
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"uid": claims.Uid,
		"sid": claims.Sid,
		"iat": time.Now().UTC().Unix(),
		"exp": claims.Expire,
		"jti": UuidNewV4().String(),
		"sig": hex.EncodeToString(sum[:]),
		"scp": claims.Scope,
	})

	block, _ := pem.Decode([]byte(claims.PrivateKey))
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	return token.SignedString(key)
}
