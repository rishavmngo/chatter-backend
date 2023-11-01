package entities

import (
	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("SecretYouShouldHide")

func GenerateToken() (string, error) {
	token := jwt.New(jwt.SigningMethodEdDSA)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
