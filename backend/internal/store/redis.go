package store

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var Client *redis.Client

func InitRedis() {
	redisAddr := os.Getenv("REDIS")
	Client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	if err := Client.Ping(Ctx).Err(); err != nil {
		panic(err)
	}
}
