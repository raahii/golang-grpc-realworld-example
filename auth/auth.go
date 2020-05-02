package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

var JWTSecret = []byte("!!SECRET!!")

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(id uint) (string, error) {
	return generateToken(id, time.Now())
}

func GenerateTokenWithTime(id uint, t time.Time) (string, error) {
	// for test
	return generateToken(id, t)
}

func generateToken(id uint, now time.Time) (string, error) {
	claims := &Claims{
		id,
		jwt.StandardClaims{
			ExpiresAt: now.Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}

	return t, nil
}

func GetUserID(ctx context.Context) (uint, error) {
	tokenString, err := grpc_auth.AuthFromMD(ctx, "Token")
	if err != nil {
		return 0, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
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

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, errors.New("invalid token: cannot map token to claims")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return 0, errors.New("token expired")
	}

	return claims.UserID, nil
}
