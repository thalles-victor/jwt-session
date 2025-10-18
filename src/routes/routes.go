package routes

import (
	domain_auth "jwt-session/src/domain/auth"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	v1 := app.Group("v1")

	//===============
	// (v1) Auth
	//================
	v1Auth := v1.Group("auth")
	v1Auth.Post("/sign-in", domain_auth.SignIn)
	v1Auth.Post("/sign-up", domain_auth.SinUp)
}
