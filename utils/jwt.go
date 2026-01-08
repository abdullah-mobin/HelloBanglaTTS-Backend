package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPayload struct {
	UserID string `json:"user_id"`
}

func GenerateAccessToken(payload TokenPayload) (any, error) {

	duration, err := time.ParseDuration(os.Getenv("JWT_EXPIRATION_TIME"))
	if err != nil {
		return nil, err
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": payload.UserID,
			"exp":     jwt.NewNumericDate(Now().Add(duration)),
		}).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	return token, nil
}

func GenerateRefreshToken(payload TokenPayload) (any, error) {

	duration, err := time.ParseDuration(os.Getenv("JWT_REFRESH_EXPIRATION_TIME"))
	if err != nil {
		return nil, err
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": payload.UserID,
			"exp":     jwt.NewNumericDate(Now().Add(duration)),
		}).SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return token, nil
}

func ValidateToken(tokenString string, secret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token claims")
	}
	userID := claims["user_id"].(string)
	return userID, nil
}
