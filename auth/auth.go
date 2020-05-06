package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// GenerateToken generates a new token
func GenerateToken(id uint) (string, error) {
	return generateToken(id, time.Now())
}

// GenerateTokenWithTime generates a new token with expired date computed withspecified time
func GenerateTokenWithTime(id uint, t time.Time) (string, error) {
	// for test
	return generateToken(id, t)
}

func generateToken(id uint, now time.Time) (string, error) {
	claims := &claims{
		id,
		jwt.StandardClaims{
			ExpiresAt: now.Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return t, nil
}

// GetUserID gets user id string from request context
func GetUserID(ctx context.Context) (uint, error) {
	tokenString, err := grpc_auth.AuthFromMD(ctx, "Token")
	if err != nil {
		return 0, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if !token.Valid {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return 0, errors.New("invalid token: it's not even a token")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return 0, errors.New("token expired")
			} else {
				return 0, fmt.Errorf("invalid token: couldn't handle this token; %w", err)
			}
		} else {
			return 0, fmt.Errorf("invalid token: couldn't handle this token; %w", err)
		}
	}

	c, ok := token.Claims.(*claims)
	if !ok {
		return 0, errors.New("invalid token: cannot map token to claims")
	}

	if c.ExpiresAt < time.Now().Unix() {
		return 0, errors.New("token expired")
	}

	return c.UserID, nil
}
