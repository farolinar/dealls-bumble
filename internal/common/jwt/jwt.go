package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/farolinar/dealls-bumble/config"
	"github.com/golang-jwt/jwt/v5"
)

var (
	// key = []byte(os.Getenv("JWT_SECRET"))

	ErrUnknownClaims = errors.New("unknown claims type")
	ErrTokenInvalid  = errors.New("invalid token")
)

func Sign(ttl time.Duration, subject string) (string, error) {
	cfg := config.GetConfig()

	now := time.Now()
	expiry := now.Add(ttl)
	t := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiry),
			Subject:   subject,
		},
	)
	return t.SignedString([]byte(cfg.App.JWTSecret))
}

func VerifyAndGetSubject(tokenString string) (string, error) {
	cfg := config.GetConfig()

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(cfg.App.JWTSecret), nil
	})
	if err != nil {
		return "", err
	}

	// Checking token validity
	if !token.Valid {
		return "", ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		return claims.Subject, nil
	} else {
		return "", ErrUnknownClaims
	}
}
