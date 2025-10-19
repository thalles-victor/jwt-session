package jwt

import (
	"context"
	"fmt"
	"jwt-session/src/utilities/config"
	"jwt-session/src/utilities/database"
	"jwt-session/src/utilities/logger"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
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

type GenerateJwtSessionProps struct {
	UserId  string
	Browser *string
	IP      *string
}

func GenerateJwtSession(p GenerateJwtSessionProps) (sessionId string, jwt string, err error) {
	logger.Info.Printf("stat the generation of session to user with id %s with browser %s and ip %s \n", p.UserId, *p.Browser, *p.IP)

	client := database.GetRedisClient()
	ctx := context.Background()

	sessionId = fmt.Sprintf("session-%s", uuid.New().String())
	logger.Info.Printf("create session id: %s \n", sessionId)

	hashFields := []string{
		"user_id", p.UserId,
		"created_at", time.Now().String(),
	}

	if p.Browser != nil {
		hashFields = append(hashFields, "browser", *p.Browser)
	}

	if p.IP != nil {
		hashFields = append(hashFields, "ip", *p.IP)
	}

	if _, err = client.HSet(ctx, sessionId, hashFields).Result(); err != nil {
		return
	}

	jwt, err = GenerateJwt(sessionId)
	if err != nil {
		return
	}

	logger.Info.Printf("create session successful. sessionId: %s", sessionId)

	return
}
