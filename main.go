package main

import (
	"fmt"
	"jwt-session/src/routes"
	"jwt-session/src/utilities/config"
	"jwt-session/src/utilities/database"
	"jwt-session/src/utilities/logger"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	logger.Init()
	config.LoadEnv()

	app := fiber.New()

	// 🔹 Servir a documentação HTML diretamente no endpoint /
	app.Static("/", "./api-doc") // <-- serve index.html da pasta api-doc

	logger.Info.Println("connecting in the database")
	db, err := database.Connect()
	if err != nil {
		logger.Error.Printf("error when connect in the database. error %s \n", err.Error())
		log.Fatalf("error when connect in the database %s", err.Error())
	}
	db.Close()

	database.GetRedisClient()

	logger.Info.Println("setup routes")
	routes.Setup(app)

	port := ":8080"
	fmt.Printf("Start the server in http://localhost%s \n", port)
	log.Fatal(app.Listen(port))
}
