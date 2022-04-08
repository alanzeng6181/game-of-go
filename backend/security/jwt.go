package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

//TODO secure this
var JWTSigningKey []byte = []byte("TempSecreteKey")

type Claims struct {
	Username string
	jwt.StandardClaims
}

func GetToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Username: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	})
	return token.SignedString(JWTSigningKey)
}

func GetUserId(tokenStr string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("expected HS356 signing method")
		}
		return JWTSigningKey, nil
	})
	if err != nil {
		return "", errors.New("error parsing token")
	}
	if token.Valid {
		return claims.Username, nil
	}
	return "", errors.New("error parsing token claims")
}
