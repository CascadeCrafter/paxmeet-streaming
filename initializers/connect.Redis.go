package initializers

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         context.Context
)

func ConnectRedis(config *Config) *redis.Client {
	Ctx = context.TODO()

	RedisClient = redis.NewClient(&redis.Options{
		Addr: config.Redis.Url,
		DB:   2, // Using database 2
	})

	if _, err := RedisClient.Ping(Ctx).Result(); err != nil {
		panic(err)
	}

	err := RedisClient.Set(Ctx, "statusHealth", "streaming proxy server is online", 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Redis client connected successfully...")

	return RedisClient
}
