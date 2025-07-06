package repository

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	Database *gorm.DB
	Redis    *redis.Client
}

func NewPaymentRepository(db *gorm.DB, redisClient *redis.Client) *PaymentRepository {
	return &PaymentRepository{
		Database: db,
		Redis:    redisClient,
	}
}
