package domain_auth

import (
	"database/sql"
	"fmt"

	"jwt-session/src/models"
	"jwt-session/src/repositories"
	"jwt-session/src/utilities/database"
	"jwt-session/src/utilities/logger"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func SinUp(c *fiber.Ctx) error {
	var body SignUpDto

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
			logger.Error.Printf("errro when get user from database. error: %s \n", err.Error())
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "erro when get user from database",
				"error":   err.Error(),
			})
		}
	}

	if err == nil {
		logger.Warn.Printf("user with email %s already registered.\n", body.Email)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("user with email %s already registered", body.Email),
		})
	}

	logger.Info.Printf("create a new user\n")

	user = models.User{
		ID:        uuid.New().String(),
		Name:      body.Name,
		Email:     body.Email,
		Password:  body.Password,
		CreatedAt: time.Now(),
	}

	userCreated, err := userRepository.Create(&user)
	if err != nil {
		logger.Error.Printf("error when create user with email %s \n", user.Email)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("internal server error when create a new user with email %s already registered", body.Email),
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "usuário cadastrado com sucesso",
		"data":    body,
		"user":    userCreated,
	})
}

func SignIn(c *fiber.Ctx) error {
	var body SignInDto

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "usuário logado com sucesso",
		"data":    body,
	})
}
