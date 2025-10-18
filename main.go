package main

import (
	"jwt-session/src/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	routes.Setup(app)

	log.Fatal((app.Listen(":8080")))
}
