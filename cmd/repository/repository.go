package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"payment_service/models"
)

type PaymentRepository interface {
	UpdateStatus(ctx context.Context, tx *gorm.DB, orderId int64, status int64) error
	WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	SavePayment(ctx context.Context, model models.Payment) error
	IsAlreadyPaid(ctx context.Context, orderID int64) (bool, error)
	GetPaymentAmountByOrderID(ctx context.Context, orderID int64) (float64, error)
	SavePaymentAnomaly(ctx context.Context, param models.PaymentAnomaly) error
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
