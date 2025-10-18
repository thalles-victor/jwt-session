package jwt

import (
	"jwt-session/src/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJwt(sub string) (string, error) {
	claims := jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(time.Hour * 1).Unix(), // expires in 1h
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWT_SEC_KEY)
}

func ParseJWT(tokenAsString string) (string, error) {
	token, err := jwt.Parse(tokenAsString, func(t *jwt.Token) (interface{}, error) {
		return config.JWT_SEC_KEY, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub := claims["user"].(string)
		return sub, nil
	}

	return "", err
}
