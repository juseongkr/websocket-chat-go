package token

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
	"os"
)

var jwtSecret = os.Getenv("SECRET")

var ErrInvalidJwtToken = errors.New("invalid jwt Token")

type customClaims struct {
	jwt.StandardClaims
	UserId int `json:"uid"`
}

func CreateNewToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
		},
	})

	return token.SignedString([]byte(jwtSecret))
}

func Parse(token string) (userId int, err error) {
	parsed, err := jwt.ParseWithClaims(token, &customClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

	if err != nil {
		return -1, err
	}

	if !parsed.Valid {
		return -1, ErrInvalidJwtToken
	}

	if c, ok := parsed.Claims.(*customClaims); ok {
		return c.UserId, nil
	}

	return -1, ErrInvalidJwtToken
}
