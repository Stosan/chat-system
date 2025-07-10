package database

import (
	"chatsystem/internal/config"
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.RedisAddress,
		Username: config.AppConfig.RedisUsername,
		Password: config.AppConfig.RedisPassword,
		DB:       0,
	})
	ctx := context.Background()
	status := rdb.Ping(ctx)
	if status.Err() != nil {
		return nil, fmt.Errorf("failed to ping Redis: %v", status.Err())
	}
	log.Println("☁️💾 \033[1;32mRedis Database ::Connected\033[0m")
	return rdb, nil
}
