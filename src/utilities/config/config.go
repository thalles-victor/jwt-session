package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DATABASE_CONNECTION = ""
	LOKI_CONNECTION     = ""
	JWT_SEC_KEY         []byte
)

func LoadEnv() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal("error when loading .env file. Check if file exist")
	}

	DATABASE_CONNECTION = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	LOKI_CONNECTION = os.Getenv("LOKI_CONNECTION")
	if LOKI_CONNECTION == "" {
		fmt.Println("LOKI_CONNECTION variable not passed in env file")
	}

	JWT_SEC_KEY = []byte(os.Getenv("JWT_SEC_KEY"))
}
