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
	UpdateStatus(ctx context.Context, tx *gorm.DB, orderId int64, status string) error

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

	// UpdatePendingPaymentRequest updates a payment request as pending for a given payment request ID.
	// ctx: Context for managing request-scoped values.
	// paymentRequestID: ID of the payment request to update.
	UpdatePendingPaymentRequest(ctx context.Context, paymentRequestID int64) error

	// GetFailedPaymentRequest retrieves all failed payment requests.
	// ctx: Context for managing request-scoped values.
	// requests: Pointer to a slice of PaymentRequests models to populate with failed requests.
	GetFailedPaymentRequest(ctx context.Context, requests *[]models.PaymentRequests) error

	// GetPaymentInfoByOrderID retrieves payment information for a given order ID.
	// ctx: Context for managing request-scoped values.
	// orderID: ID of the order.
	// Returns: A Payment model containing payment details and an error if any issues occur.
	GetPaymentInfoByOrderID(ctx context.Context, orderID int64) (models.Payment, error)

	// MarkExpiredPayments marks payments as expired for a given payment ID.
	// ctx: Context for managing request-scoped values.
	// paymentId: ID of the payment to mark as expired.
	MarkExpiredPayments(ctx context.Context, paymentId int64) error

	// GetExpiredPendingPayments retrieves all expired payments.
	// ctx: Context for managing request-scoped values.
	// Returns: A slice of Payment models containing expired payments and an error if any issues occur.
	GetExpiredPendingPayments(ctx context.Context) ([]models.Payment, error)

	// InsertAuditLog inserts an audit log entry into the database.
	// ctx: Context for managing request-scoped values.
	// param: PaymentAuditLog model containing audit log details.
	InsertAuditLog(ctx context.Context, param models.PaymentAuditLog) error
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
