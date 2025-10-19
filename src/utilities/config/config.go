package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	DATABASE_CONNECTION     string
	GRAFANA_LOKI_CONNECTION string
	JWT_SEC_KEY             []byte
	MAIL_HOST               string
	MAIL_USER               string
	MAIL_PASS               string
	MAIL_PORT               int

	REDIS_ADDR string
	REDIS_USER string
	REDIS_pass string
	REDIS_DB   int
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

	loadMailCredentialsWithPanicOnError()
	loadRedisEnvWithPanicOnError()

	GRAFANA_LOKI_CONNECTION = os.Getenv("LOKI_CONNECTION")
	if GRAFANA_LOKI_CONNECTION == "" {
		fmt.Println("GRAFANA_LOKI_CONNECTION variable not passed in env file")
	}

	JWT_SEC_KEY = []byte(os.Getenv("JWT_SEC_KEY"))
}

func loadMailCredentialsWithPanicOnError() error {
	MAIL_HOST = os.Getenv("MAIL_HOST")
	if MAIL_HOST == "" {
		log.Fatal("MAIL_HOST not declared in environment variable")
	}

	MAIL_USER = os.Getenv("MAIL_USER")
	if MAIL_HOST == "" {
		log.Fatal("MAIL_USER not declared in environment variable")
	}

	MAIL_PASS = os.Getenv("MAIL_PASS")
	if MAIL_HOST == "" {
		log.Fatal("MAIL_PASS not declared in environment variable")
	}

	mailPort, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		log.Fatal("MAIL_PASS no a number declared")
	}
	MAIL_PORT = mailPort

	return nil
}

func loadRedisEnvWithPanicOnError() {
	addr := os.Getenv("REDIS_ADDR")
	user := os.Getenv("REDIS_USER")
	pass := os.Getenv("REDIS_PASS")

	if addr == "" || user == "" || pass == "" {
		log.Fatal("estão faltando váriaveis do redis, checa ai meu parceiro")
	}

	REDIS_ADDR = addr
	REDIS_USER = user
	REDIS_pass = pass
	REDIS_DB = 0
}
