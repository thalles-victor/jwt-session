package main

import (
	"jwt-session/src/logger"
	"jwt-session/src/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	logger.Init()

	routes.Setup(app)

	port := ":8080"
	logger.Info.Printf("Start the server in http://localhost%s", port)
	log.Fatal((app.Listen(port)))
}
