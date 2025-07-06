package resource

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"payment_service/config"
)

func InitRedis(cfg *config.Config) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
	})

	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect with reds %s", err)
	}
	log.Print("Connected with redis")
	return redisClient
}
