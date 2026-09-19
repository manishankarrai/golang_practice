package db

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

// RedisClient is the shared Redis client used for vehicle-location caching.
var RedisClient *redis.Client

// ConnectRedis initializes Redis using REDIS_URL or localhost:6379 and verifies connectivity.
func ConnectRedis() *redis.Client {
	options, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil || options == nil {
		options = &redis.Options{Addr: "localhost:6379"}
	}
	RedisClient = redis.NewClient(options)
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	return RedisClient
}
