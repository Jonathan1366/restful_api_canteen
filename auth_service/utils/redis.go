package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	url := os.Getenv("REDIS_URL")
	opt, err := redis.ParseURL(url)
	if err != nil {
		fmt.Println("Failed to parse Redis URL: ", err)
		return
	}
	RedisClient = redis.NewClient(opt)

	_, err = RedisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Redis connection failed:", err)

	} else{
		fmt.Println("Redis connected successfully")
	}
}
