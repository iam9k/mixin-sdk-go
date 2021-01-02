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
	uid        string
	sid        string
	privateKey string
	method     string
	uri        string
	body       string
	scope      string
	expire     int64
}

func SignAuthenticationToken(claims SignClaims) (string, error) {
	sum := sha256.Sum256([]byte(claims.method + claims.uri + claims.body))
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"uid": claims.uid,
		"sid": claims.sid,
		"iat": time.Now().UTC().Unix(),
		"exp": claims.expire,
		"jti": UuidNewV4().String(),
		"sig": hex.EncodeToString(sum[:]),
		"scp": claims.scope,
	})

	block, _ := pem.Decode([]byte(claims.privateKey))
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	return token.SignedString(key)
}
