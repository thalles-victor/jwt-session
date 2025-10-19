package database

import (
	"context"
	"fmt"
	"jwt-session/src/utilities/config"
	"log"

	"github.com/redis/go-redis/v9"
)

func GetRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.REDIS_ADDR,
		Username: config.REDIS_USER,
		Password: config.REDIS_pass,
		DB:       config.REDIS_DB,
	})

	ctx := context.Background()
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Erro ao conectar no Redis: %v", err)
	}
	fmt.Println("Conectado ao Redis:", pong)

	return client
}
