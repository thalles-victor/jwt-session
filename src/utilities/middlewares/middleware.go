package middlewares

import (
	"context"
	"jwt-session/src/utilities/database"
	"jwt-session/src/utilities/jwt"
	"jwt-session/src/utilities/logger"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JwtMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	logger.Info.Println("extract auth header")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing authorization header",
		})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token format",
		})
	}

	token := parts[1]

	sub, err := jwt.ParseJWT(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid or expired token",
		})
	}

	c.Locals("userId", sub)

	return c.Next()
}

func JwtSessionMiddleware(c *fiber.Ctx) error {
	logger.Info.Println("extract auth header")
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing authorization header",
		})
	}

	logger.Info.Println("extract token from authorization")
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	logger.Info.Printf("try extract sessionId and check if token is valid")
	sessionId, err := jwt.ParseJWT(tokenStr)
	if err != nil || sessionId == "" {
		logger.Warn.Print("token invalid")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "token de autenticação inválido",
			"error":   err.Error(),
		})
	}
	logger.Info.Printf("extracted session with id: %s ", sessionId)

	logger.Info.Println("connect in redis client")
	client := database.GetRedisClient()
	ctx := context.Background()

	logger.Info.Printf("check if session exist")
	sessionData, err := client.HGetAll(ctx, sessionId).Result()
	if err != nil || len(sessionData) == 0 {
		logger.Warn.Printf("session not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "sessão inválida ou expirada",
			"error":   err.Error(),
		})
	}

	logger.Info.Println("extract user_id from session")
	userId, ok := sessionData["user_id"]
	if !ok || userId == "" {
		logger.Info.Println("user_id not found in the session")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "usuário não encontrado na sessão",
			"error":   err.Error(),
		})
	}

	c.Locals("userId", userId)

	logger.Info.Printf("user with id %s authorized", userId)
	return c.Next()
}
