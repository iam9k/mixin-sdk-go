package mixin

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func SignAuthenticationToken(uid, sid, privateKey, method, uri, body string) (string, error) {
	expire := time.Now().UTC().Add(time.Hour * 24 * 30 * 3)
	sum := sha256.Sum256([]byte(method + uri + body))
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"uid": uid,
		"sid": sid,
		"iat": time.Now().UTC().Unix(),
		"exp": expire.Unix(),
		"jti": UuidNewV4().String(),
		"sig": hex.EncodeToString(sum[:]),
		"scp": "FULL",
	})

	block, _ := pem.Decode([]byte(privateKey))
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	return token.SignedString(key)
}
