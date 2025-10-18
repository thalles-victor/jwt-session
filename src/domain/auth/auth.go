package domain_auth

import "github.com/gofiber/fiber/v2"

func SinUp(c *fiber.Ctx) error {
	var body SignUpDto

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "usuário cadastrado com sucesso",
		"data":    body,
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
