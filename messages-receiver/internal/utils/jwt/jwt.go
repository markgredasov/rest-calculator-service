package utils_jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const keyString = "JWT_SECRET_KEY"

type CustomClaims struct {
	*jwt.RegisteredClaims
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

func GenerateToken(userID, role string) (string, error) {
	secretKey := os.Getenv(keyString)
	if secretKey == "" {
		return "", fmt.Errorf("no JWT_SECRET_KEY provided")
	}

	claims := CustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour).UTC()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	stringToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("get jwt signed string: %w", err)
	}

	return stringToken, nil
}

func ExtractJWTClaims(tokenString string) (*CustomClaims, error) {
	secretKey := os.Getenv(keyString)
	if secretKey == "" {
		return nil, fmt.Errorf("no JWT_SECRET_KEY provided")
	}

	var claims CustomClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("JWT expired: %w", err)
		}

		return nil, fmt.Errorf("parse with claims error: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token invalid")
	}

	return &claims, nil
}
