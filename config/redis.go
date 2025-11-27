package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(
		&redis.Options{
			Addr:     "127.0.0.1:6379",
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
