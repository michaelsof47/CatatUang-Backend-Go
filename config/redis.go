package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Context = context.Background()

func InitRedis() *redis.Client {
	redisdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if err := redisdb.Ping(Context).Err(); err != nil {
		log.Fatal("Redis Connection Failed:", err)
	}

	log.Println("Redis Connected")
	return redisdb
}
