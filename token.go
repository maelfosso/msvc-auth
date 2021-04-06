package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var _JWTPrivateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))
var _JWTPublicKey = []byte(os.Getenv("JWT_PUBLIC_KEY"))

type CustomClaims struct {
	*jwt.StandardClaims
	User
}

func GetJWT(user User) (string, error) {
	signBytes, err := ioutil.ReadFile(string(_JWTPrivateKey))
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

func VerifyJWT(token string) (interface{}, error) {
	verifyBytes, err := ioutil.ReadFile(string(_JWTPublicKey))
	if err != nil {
		return "", err
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return "", fmt.Errorf("Validate: parse key: %w", err)
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected method: %s", jwtToken.Header["alg"])
		}

		return verifyKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Validate: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("Validate: invalid")
	}

	return claims["dat"], nil
}
