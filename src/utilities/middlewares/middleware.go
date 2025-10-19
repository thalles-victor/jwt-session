package middlewares

import (
	"database/sql"
	"jwt-session/src/models"
	"jwt-session/src/repositories"
	"jwt-session/src/utilities/database"
	"jwt-session/src/utilities/date"
	"jwt-session/src/utilities/jwt"
	"jwt-session/src/utilities/logger"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JwtMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

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
	authHeader := c.Get("Authorization")

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

	logger.Info.Println("check if session exist")

	db, err := database.Connect()
	if err != nil {
		logger.Error.Printf("error when connect in the databse. err: %s \n", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "erro when connect in the database",
			"error":   err.Error(),
		})
	}
	defer db.Close()
	sessionRepository := repositories.NewSessionRepository(db)

	var session models.Session
	if err = sessionRepository.GetByID(sub, &session); err != nil && err != sql.ErrNoRows {
		logger.Error.Printf("internal server error when get session by id %s. error: %s", sub, err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "erro interno no servidor ao tentar buscar a sessão",
			"error":   err.Error(),
		})
	}

	if err == sql.ErrNoRows {
		logger.Warn.Printf("session with id %s no found", sub)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "sessão não encontrada",
		})
	}

	if !date.IsNotExpired(session.ExpiresAt) {
		logger.Warn.Printf("session already expired at %s", session.ExpiresAt)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "sessão inválida ou expirada",
		})
	}

	c.Locals("userId", session.UserID)
	c.Locals("sessionId", session.ID)

	return c.Next()
}
