package domain_auth

import (
	"database/sql"
	"fmt"
	"jwt-session/src/models"
	"jwt-session/src/repositories"
	"jwt-session/src/utilities/database"
	"jwt-session/src/utilities/logger"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func GetAllSessionsFromUser(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)

	db, err := database.Connect()
	if err != nil {
		logger.Error.Printf("error when connect in the databse. err: %s \n", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "erro when connect in the database",
			"error":   err.Error(),
		})
	}
	defer db.Close()

	userRepository := repositories.NewUserRepository(db)

	var user models.User
	if err = userRepository.GetById(userId, &user); err != nil {
		if err == sql.ErrNoRows {
			logger.Warn.Printf("user unregistered")
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": "usuário não cadastrado",
			})
		}

		logger.Warn.Printf("user unregistered")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "erro interno no servidor ao buscar usuário",
			"error":   err.Error(),
		})
	}

	sessionRepository := repositories.NewSessionRepository(db)
	var sessions []models.Session
	if err = sessionRepository.GetAllByUserID(userId, &sessions); err != nil {
		logger.Error.Printf("internal server error when get all sessions from user. error: %s ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "error interno no servidor ao buscar todas as sessões",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":  "sessões encontradas com successo",
		"sessions": sessions,
	})
}

func DeleteSessionFromUser(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	sessionId := c.Params("sessionId")

	db, err := database.Connect()
	if err != nil {
		logger.Error.Printf("error when connect in the databse. err: %s \n", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "erro when connect in the database",
			"error":   err.Error(),
		})
	}
	defer db.Close()

	userRepository := repositories.NewUserRepository(db)

	var user models.User
	if err = userRepository.GetById(userId, &user); err != nil {
		if err == sql.ErrNoRows {
			logger.Warn.Printf("user unregistered")
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": "usuário não cadastrado",
			})
		}

		logger.Warn.Printf("user unregistered")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "erro interno no servidor ao buscar usuário",
			"error":   err.Error(),
		})
	}

	sessionRepository := repositories.NewSessionRepository(db)

	var session models.Session
	if err = sessionRepository.GetByID(sessionId, &session); err != nil {
		if err != sql.ErrNoRows {
			logger.Error.Printf("internal server when get session from database. error %s", err.Error())
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "error interno no servidor ao buscar sessão",
				"error":   err.Error(),
			})
		}

		logger.Warn.Printf("session not found. error %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("sessão com id: %s não encontrada", sessionId),
			"error":   err.Error(),
		})
	}

	if session.UserID != userId {
		logger.Warn.Printf("the session with id %s belong another user", sessionId)
		return c.Status(http.StatusNotAcceptable).JSON(fiber.Map{
			"message": fmt.Sprintf("a sessão com id %s pertence a outro usuário", sessionId),
		})
	}

	if err = sessionRepository.DeleteByID(sessionId); err != nil {
		logger.Error.Printf("internal server error when delete session of user. error: %s ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "error interno no servidor ao buscar todas as sessões",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "sessão deletada com sucesso",
	})
}
