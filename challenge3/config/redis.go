package config

import (
	"context"
	// "fmt"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
func InitRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		// Addr:     "localhost:6379",
		Addr:     "redis:6379",
		Password: "", // no password
		DB:       0,  // use default DB
		Protocol: 2,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	RedisClient = rdb
	return rdb, nil
}
