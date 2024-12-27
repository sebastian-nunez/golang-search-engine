package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sebastian-nunez/golang-search-engine/config"
)

type AuthClaims struct {
	ID      string `json:"id"`
	User    string `json:"user"`
	IsAdmin bool   `json:"isAdmin"`
	jwt.RegisteredClaims
}

func CreateNewAuthToken(id string, email string, isAdmin bool) (string, error) {
	claims := AuthClaims{
		ID:      id,
		User:    email,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "sebastian-nunez.com",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.Envs.SecretKey))
	if err != nil {
		return "", fmt.Errorf("unable to sign JWT token: %s", err)
	}

	return tokenStr, nil
}

func ValidateJWT(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Envs.SecretKey), nil
	})
}
