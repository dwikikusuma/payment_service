package repository

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type OrderRepository struct {
	Database *gorm.DB
	Redis    *redis.Client
}

func NewOrderRepository(db *gorm.DB, redisClient *redis.Client) *OrderRepository {
	return &OrderRepository{
		Database: db,
		Redis:    redisClient,
	}
}
