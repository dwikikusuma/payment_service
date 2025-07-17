package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"payment_service/models"
)

// PaymentRepository defines the interface for payment-related database operations.
type PaymentRepository interface {
	// UpdateStatus updates the status of a payment for a given order ID.
	// ctx: Context for managing request-scoped values.
	// tx: Database transaction object.
	// orderId: ID of the order.
	// status: New status to be updated.
	UpdateStatus(ctx context.Context, tx *gorm.DB, orderId int64, status int64) error

	// WithTransaction executes a function within a database transaction.
	// ctx: Context for managing request-scoped values.
	// fn: Function to execute within the transaction.
	WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error

	// SavePayment saves a payment record to the database.
	// ctx: Context for managing request-scoped values.
	// model: Payment model containing payment details.
	SavePayment(ctx context.Context, model models.Payment) error

	// IsAlreadyPaid checks if a payment has already been made for a given order ID.
	// ctx: Context for managing request-scoped values.
	// orderID: ID of the order.
	IsAlreadyPaid(ctx context.Context, orderID int64) (bool, error)

	// GetPaymentAmountByOrderID retrieves the payment amount for a given order ID.
	// ctx: Context for managing request-scoped values.
	// orderID: ID of the order.
	GetPaymentAmountByOrderID(ctx context.Context, orderID int64) (float64, error)

	// SavePaymentAnomaly saves a record of a payment anomaly to the database.
	// ctx: Context for managing request-scoped values.
	// param: Payment anomaly model containing anomaly details.
	SavePaymentAnomaly(ctx context.Context, param models.PaymentAnomaly) error

	// SaveFailedPublishEvent saves a record of a failed event to the database.
	// ctx: Context for managing request-scoped values.
	// param: Failed event model containing event details.
	SaveFailedPublishEvent(ctx context.Context, param models.FailedEvents) error

	// GetPendingPayment retrieves all pending payments from the database.
	// ctx: Context for managing request-scoped values.
	GetPendingPayment(ctx context.Context) ([]models.Payment, error)

	// UpdateSuccessPaymentRequest updates a payment request as successful for a given order ID.
	// ctx: Context for managing request-scoped values.
	// orderID: ID of the order.
	UpdateSuccessPaymentRequest(ctx context.Context, orderID int64) error

	// UpdateFailedPaymentRequest updates a payment request as failed for a given order ID.
	// ctx: Context for managing request-scoped values.
	// orderID: ID of the order.
	// notes: Notes describing the failure.
	UpdateFailedPaymentRequest(ctx context.Context, orderID int64, notes string) error

	// GetPendingPaymentRequest retrieves all pending payment requests.
	// ctx: Context for managing request-scoped values.
	// paymentRequests: Pointer to a slice of PaymentRequests models to populate.
	GetPendingPaymentRequest(ctx context.Context, paymentRequests *[]models.PaymentRequests) error

	// SavePaymentRequest saves a payment request record to the database.
	// ctx: Context for managing request-scoped values.
	// param: Payment request model containing request details.
	SavePaymentRequest(ctx context.Context, param models.PaymentRequests) error
}

// paymentRepository is the implementation of the PaymentRepository interface.
type paymentRepository struct {
	Database *gorm.DB      // Database connection object.
	Redis    *redis.Client // Redis client for caching or other operations.
}

// NewPaymentRepository creates a new instance of paymentRepository.
// db: Database connection object.
// redisClient: Redis client instance.
func NewPaymentRepository(db *gorm.DB, redisClient *redis.Client) PaymentRepository {
	return &paymentRepository{
		Database: db,
		Redis:    redisClient,
	}
}
