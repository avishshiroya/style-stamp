package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthDto struct {
	Username string    `json:"username"`
	ID       uuid.UUID `json:"id"`
	Jti      uuid.UUID `json:"jti"`
}

func CreateAccessToken(auth AuthDto) (string, error) {
	var secretKey = []byte("83f138c1-801b-4f27-bcd6-ee0dca60d349")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":       auth.ID,
			"username": auth.Username,
			"jti":      auth.Jti,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		},
	)
	tokenstring, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenstring, nil
}
func CreateRefreshToken(auth AuthDto) (string, error) {
	var secretKey = []byte("02608933-734C-45F7-B0F3-5CD288E36774")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":       auth.ID,
			"username": auth.Username,
			"jti":      auth.Jti,
			"exp":      time.Now().Add(time.Hour * 24 * 2).Unix(),
		},
	)
	tokenstring, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenstring, nil
}

func VerifyAccessToken(s string) (jwt.MapClaims, error) {
	var secretKey = []byte("83f138c1-801b-4f27-bcd6-ee0dca60d349")
	tokens, err := jwt.Parse(s, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := tokens.Claims.(jwt.MapClaims); ok && tokens.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
func VerifyRefreshToken(s string) (jwt.MapClaims, error) {
	var secretKey = []byte("02608933-734C-45F7-B0F3-5CD288E36774")
	tokens, err := jwt.Parse(s, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := tokens.Claims.(jwt.MapClaims); ok && tokens.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
