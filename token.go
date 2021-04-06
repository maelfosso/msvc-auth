package main

import (
	"io/ioutil"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var _JWTSignInKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type CustomClaims struct {
	*jwt.StandardClaims
	User
}

func GetJWT(user User) (string, error) {
	signBytes, err := ioutil.ReadFile(string(_JWTSignInKey))
	if err != nil {
		return "", err
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()
	ttl := time.Hour

	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims = &CustomClaims{
		&jwt.StandardClaims{
			ExpiresAt: now.Add(ttl).Unix(),
			IssuedAt:  now.Unix(),
		},
		user,
	}

	return token.SignedString(signKey)
}

// func VerifyJWT(string) (bool, error) {

// }
