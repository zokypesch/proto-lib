package core

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Jwt struct info for JWT
type Jwt struct {
	ID    int
	Email string
	Name  string
}

var secret = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

// GenerateToken for generate a token
func GenerateToken(info Jwt) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    info.ID,
		"email": info.Email,
		"name":  info.Name,
		"nbf":   time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(secret)
}

// ExtractToken for extract a token
func ExtractToken(tokenString string) (Jwt, error) {
	token, err := GetToken(tokenString)

	if err != nil {
		return Jwt{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return Jwt{}, fmt.Errorf("token not valid")
	}

	return Jwt{
		ID:    int(claims["id"].(float64)),
		Email: claims["email"].(string),
		Name:  claims["name"].(string),
	}, nil

}
