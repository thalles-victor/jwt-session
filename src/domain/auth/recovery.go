package domain_auth

import (
	"database/sql"
	"fmt"
	"jwt-session/src/models"
	"jwt-session/src/repositories"
	"jwt-session/src/utilities/code"
	"jwt-session/src/utilities/database"
	"jwt-session/src/utilities/date"
	"jwt-session/src/utilities/logger"
	"jwt-session/src/utilities/mail"
	"net/http"
	"time"

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
			logger.Error.Printf("user with email %s un registered. \n", email)
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

	expiresAt, err := date.GenerateFutureDate(5, "minutes")
	if err != nil {
		logger.Error.Printf("erro ao gerar data futura %s \n", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "error interno no servidor ao gerar data futura para o codigo",
			"error":   err.Error(),
		})
	}

	recoveryRepository := repositories.NewRecoveryRepository(db)

	logger.Info.Printf("check if recovery exist by email %s \n", email)
	var recovery models.Recovery
	err = recoveryRepository.GetByEmail(email, &recovery)
	if err != nil && err != sql.ErrNoRows {
		logger.Error.Printf("internal server error when get recovery data")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "error interno no servidor ao buscar os dados de recupearção",
			"error":   err.Error(),
		})
	}

	logger.Info.Println("save recovery data in the database")
	if err == sql.ErrNoRows {
		logger.Info.Println("recovery data not found, then create")

		if _, err := recoveryRepository.Create(&models.Recovery{
			ID:        0,
			UserID:    user.ID,
			Email:     user.Email,
			Code:      recoveryCode,
			Attempts:  0,
			ExpiresAt: expiresAt,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Expired:   false,
		}); err != nil {
			logger.Error.Printf("error when save recovery data in the database. error: %s \n", err.Error())
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "erro inteno no servidor ao salvar os dados de recupearção",
				"error":   err.Error(),
			})
		}
	} else {
		logger.Info.Println("data of recovery already exist, then update to new data.")
		if err = recoveryRepository.UpdateRecovery(recovery.ID, recoveryCode, 0, expiresAt, false); err != nil {
			logger.Error.Printf("intternal server errro when update recovery table %s \n", err.Error())
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "erro ao atulizar os dados de recuperação",
				"error":   err.Error(),
			})
		}
	}

	recoveryUrl := fmt.Sprintf("https://dominio-do-front/recuperar/tocar-senha?code=%s", recoveryCode)

	go func() { mail.SendRecoveryRequestEmail(user.Name, user.Email, recoveryUrl) }()

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("o codigo de recuperacao foi enviado no email %s com sucesso", email),
	})
}

func ChangePasswordRequestRecovery(c *fiber.Ctx) error {
	var body ChangePasswordRequestRecoveryDto

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

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

	logger.Info.Printf("check if user with email %s exist", body.Email)
	var user models.User
	if err = userRepository.GetByEmail(body.Email, &user); err != nil {
		if err != sql.ErrNoRows {
			logger.Warn.Printf("user with email %s un registered. \n", body.Email)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "error interno no servidor ao buscar usuário",
				"error":   err.Error(),
			})
		}

		logger.Error.Printf("errro when get user from database. error: %s \n", err.Error())
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("usuário com email %s não foi cadastrado", body.Email),
		})
	}

	recoveryRepository := repositories.NewRecoveryRepository(db)

	var recovery models.Recovery
	if err := recoveryRepository.GetByEmail(body.Email, &recovery); err != nil {
		if err != sql.ErrNoRows {
			logger.Error.Printf("internal server error whe get recovery data: %s \n", err.Error())
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "error interno no servidor ao dados de recuperação",
				"error":   err.Error(),
			})
		}

		logger.Warn.Printf("recovery data no found: %s \n", err.Error())
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "dados de recuperação não encontrado, \n solicite uma recuperação de seha antes",
			"error":   err.Error(),
		})
	}

	if recovery.Expired {
		logger.Warn.Printf("recovery code already expires")
		return c.Status(http.StatusNotAcceptable).JSON(fiber.Map{
			"message": "codigo de recuperação já expirou por tentativa ou data, gere um novo",
		})
	}

	if recovery.Attempts > 10 {
		logger.Warn.Printf("number of attempts exceed, clear recovery data")
		if err = recoveryRepository.MarkAsExpired(recovery.ID); err != nil {
			logger.Error.Printf("internal server error when make make as expired. error %s", err.Error())
			return c.Status(http.StatusNotAcceptable).JSON(fiber.Map{
				"message": "erro interno no servidor ao expierar o token",
			})
		}
		return c.Status(http.StatusNotAcceptable).JSON(fiber.Map{
			"message": "número de tentativas excedido gere um novo codigo de recuperação",
		})
	}

	if !date.IsNotExpired(recovery.ExpiresAt) {
		if err = recoveryRepository.MarkAsExpired(recovery.ID); err != nil {
			logger.Error.Printf("internal server error when disable recovery code")
			c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"messge": "erro interno no servidor ao desativer o codigo de recuperação",
				"error":  err.Error(),
			})
		}

		logger.Warn.Printf("recovery code already expires at: %s", recovery.ExpiresAt)
		return c.Status(http.StatusNotAcceptable).JSON(fiber.Map{
			"message": "codigo de recuperação já expirou",
		})
	}

	logger.Info.Println("check if the code are same")
	if body.Code != recovery.Code {
		logger.Warn.Println("the recovery are differents. Increasing the attempts.")
		if err := recoveryRepository.IncrementAttempts(recovery.ID); err != nil {
			logger.Error.Printf("error when increments attempts in recovery table. error: %s \n", err.Error())
			return c.Status(http.StatusNotAcceptable).JSON(fiber.Map{
				"message": "erro interno no servidor ao incrimentar tentativas",
				"error":   err.Error(),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("o codigo informado é inválido, restam %d tentativas", 10-recovery.Attempts+1),
		})
	}

	logger.Info.Printf("code valid to user %s. Update password", user.Email)
	if err := userRepository.UpdatePassword(user.ID, body.NewPassword); err != nil {
		logger.Error.Printf("internal server error when update user password. error: %s \n", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "error interno no servidor ao atulizar a senha do usuário",
			"error":   err.Error(),
		})
	}

	logger.Info.Printf("deleting all sessions from use with id: %s", user.ID)
	sessionRepository := repositories.NewSessionRepository(db)
	if err = sessionRepository.DeleteByUserId(user.ID); err != nil {
		logger.Error.Printf("internal server error when delete all sessions from user. error: %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "erro interno no sevidor ao deletar todas as sessões de um usuário",
			"error":   err.Error(),
		})
	}

	logger.Info.Println("disble recovery code")
	if err = recoveryRepository.MarkAsExpired(recovery.ID); err != nil {
		logger.Error.Printf("internal server error when disable recovery code")
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"messge": "erro interno no servidor ao desativer o codigo de recuperação",
			"error":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "senha alterada com sucesso",
	})
}
