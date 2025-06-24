package store

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var Client *redis.Client

func InitRedis() {

	Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	if err := Client.Ping(Ctx).Err(); err != nil {
		panic(err)
	}
}
