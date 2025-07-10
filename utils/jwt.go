package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("your_super_secret_key")

type Claims struct {
	UserID  int    `json:"user_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int, name, email string, isAdmin bool) (string, error) {
	expirationTime := time.Now().Add(60 * 24 * time.Hour)
	claims := &Claims{
		UserID:  userID,
		Name:    name,
		Email:   email,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

func ParseJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
