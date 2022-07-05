package auth

import (
	"fmt"
	"time"

	"social-network/internal/errapp"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

func (s *service) generateToken(login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenExpired).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Username: login,
	})

	tokenStr, err := token.SignedString(s.signingKey)
	if err != nil {
		return "", fmt.Errorf("token.SignedString error: %v", err)
	}

	return tokenStr, nil
}

func (s *service) parseToken(accessToken string) error {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method %v", token.Header["alg"])
		}
		return s.signingKey, nil
	})
	if err != nil {
		if jwtErr, ok := err.(*jwt.ValidationError); ok {
			if jwtErr.Errors == jwt.ValidationErrorExpired {
				return errapp.ExpiredToken
			}
		}
		return fmt.Errorf("jwt.ParseWithClaims error: %v", err)
	}

	if _, ok := token.Claims.(*Claims); !(ok && token.Valid) {
		return errapp.InvalidToken
	}

	return nil
}
