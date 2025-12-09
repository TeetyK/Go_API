package config

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "127.0.0.1"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	RedisClient = redis.NewClient(
		&redis.Options{
			Addr:     redisAddr,
			Password: "",
			DB:       0,
		})
	pong, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("WARNING: Failed to connect to Redis: %v. Caching disabled.\n", err)
		RedisClient = nil
	} else {
		fmt.Println("Redis connected successfully:", pong)
	}
}
