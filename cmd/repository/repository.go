package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	UpdateStatus(ctx context.Context, tx *gorm.DB, orderId int64, status int64) error
	WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}

type paymentRepository struct {
	Database *gorm.DB
	Redis    *redis.Client
}

func NewPaymentRepository(db *gorm.DB, redisClient *redis.Client) PaymentRepository {
	return &paymentRepository{
		Database: db,
		Redis:    redisClient,
	}
}
