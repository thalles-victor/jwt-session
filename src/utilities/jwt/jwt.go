package jwt

import (
	"fmt"
	"jwt-session/src/utilities/config"
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
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, ok := claims["sub"].(string)
		if !ok {
			return "", fmt.Errorf("sub claim not found or invalid")
		}
		return sub, nil
	}

	return "", fmt.Errorf("invalid token")
}
