package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var Context = context.Background()

func InitRedis() *redis.Client {

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	redisdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	if err := redisdb.Ping(Context).Err(); err != nil {
		log.Fatal("Redis Connection Failed:", err)
	}

	log.Println("Redis Connected")
	return redisdb
}
