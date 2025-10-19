package routes

import (
	domain_auth "jwt-session/src/domain/auth"
	"jwt-session/src/utilities/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	v1 := app.Group("v1")

	//====================================================================================
	// (v1) Auth
	//====================================================================================
	{
		v1Auth := v1.Group("auth")
		v1Auth.Post("/sign-in", domain_auth.SignIn)
		v1Auth.Post("/sign-up", domain_auth.SinUp)

		//====================================================================================
		// (v1) Session Auth
		//====================================================================================
		v1SessionAuth := v1Auth.Group("/session")
		v1SessionAuth.Post("/sign-in", domain_auth.SignInWithSession)
		v1SessionAuth.Post("/sign-up", domain_auth.SinUpWithSession)
		v1SessionAuth.Get("/", middlewares.JwtSessionMiddleware, func(c *fiber.Ctx) error {
			userId := c.Locals("userId").(string)

			return c.JSON(fiber.Map{
				"message": "sessão válida",
				"user_id": userId,
			})
		})

		//====================================================================================
		// (v1) Recovery Auth
		//====================================================================================
		v1AuthRecovery := v1Auth.Group("/recovery")
		v1AuthRecovery.Post("/request/:email", domain_auth.RequestRecovery)
	}
}
