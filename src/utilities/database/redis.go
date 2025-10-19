package database

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func GetRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Username: "admin",
		Password: "minha_senha_admin",
		DB:       0,
	})

	ctx := context.Background()
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Erro ao conectar no Redis: %v", err)
	}
	fmt.Println("Conectado ao Redis:", pong)

	return client
}
