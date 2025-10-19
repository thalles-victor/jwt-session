package domain_auth

import (
	"database/sql"
	"fmt"
	"jwt-session/src/models"
	"jwt-session/src/repositories"
	"jwt-session/src/utilities/code"
	"jwt-session/src/utilities/database"
	"jwt-session/src/utilities/logger"
	"jwt-session/src/utilities/mail"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func RequestRecovery(c *fiber.Ctx) error {
	email := c.Params("email")

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

	logger.Info.Printf("check if user with email %s exist", email)
	var user models.User
	if err = userRepository.GetByEmail(email, &user); err != nil {
		if err != sql.ErrNoRows {
			logger.Warn.Printf("user with email %s un registered. \n", email)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "error interno no servidor ao buscar usuoario",
				"error":   err.Error(),
			})
		}

		logger.Error.Printf("errro when get user from database. error: %s \n", err.Error())
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("usuário com email %s não foi cadastrado", email),
		})
	}

	recoveryCode, err := code.GenerateRecoveryCode(120)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "error interno no servidor ao gerar o codigo de recuperação",
			"error":   err.Error(),
		})
	}

	recoveryUrl := fmt.Sprintf("https://dominio-do-front/recuperar/tocar-senha?code=%s", recoveryCode)

	go func() { mail.SendRecoveryRequestEmail(user.Name, user.Email, recoveryUrl) }()

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("o codigo de recuperacao foi enviado no email %s com sucesso", email),
	})
}

func ChangePasswordRequestRecovery() {}
