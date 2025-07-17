package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
	"payment_service/cmd/repository"
	"payment_service/infra/constant"
	"payment_service/infra/log"
	"payment_service/models"
	"time"
)

const (
	// maxRetryPublish defines the maximum number of retries for publishing events.
	maxRetryPublish = 5
)

// PaymentService defines the interface for payment-related operations.
type PaymentService interface {
	// ProcessPaymentSuccess processes a successful payment for a given order ID.
	// ctx: Context for managing request-scoped values.
	// orderId: ID of the order.
	// status: Status of the payment.
	ProcessPaymentSuccess(ctx context.Context, orderId int64, status string) error

	// CheckPaymentAmountByOrderID retrieves the payment amount for a given order ID.
	// ctx: Context for managing request-scoped values.
	// orderID: ID of the order.
	CheckPaymentAmountByOrderID(ctx context.Context, orderID int64) (float64, error)

	// SaveAnomaly saves a payment anomaly record to the database.
	// ctx: Context for managing request-scoped values.
	// param: Payment anomaly model containing anomaly details.
	SaveAnomaly(ctx context.Context, param models.PaymentAnomaly) error

	// SavePaymentRequest saves a payment request record to the database.
	// ctx: Context for managing request-scoped values.
	// param: Payment request model containing request details.
	SavePaymentRequest(ctx context.Context, param models.PaymentRequests) error
}

// paymentService is the implementation of the PaymentService interface.
type paymentService struct {
	PaymentRepository repository.PaymentRepository     // Repository for payment-related database operations.
	Publisher         repository.PaymentEventPublisher // Publisher for payment-related events.
}

// NewPaymentService creates a new instance of paymentService.
// paymentRepo: Repository for payment-related database operations.
// publisher: Publisher for payment-related events.
func NewPaymentService(paymentRepo repository.PaymentRepository, publisher repository.PaymentEventPublisher) PaymentService {
	return &paymentService{
		PaymentRepository: paymentRepo,
		Publisher:         publisher,
	}
}

// ProcessPaymentSuccess processes a successful payment for a given order ID.
// It updates the payment status, checks if the payment is already made, and publishes a success event.
// ctx: Context for managing request-scoped values.
// orderId: ID of the order.
// status: Status of the payment.
func (s *paymentService) ProcessPaymentSuccess(ctx context.Context, orderId int64, status string) error {
	statusID, err := constant.TranslateStatusByName(status)
	if err != nil {
		return err
	}

	isPaid, err := s.PaymentRepository.IsAlreadyPaid(ctx, orderId)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"order_id": orderId,
		}).Errorf("error occurred on s.PaymentRepository.IsAlreadyPaid(ctx, orderId): %s", err)
		return err
	}

	if isPaid {
		log.Logger.WithFields(logrus.Fields{
			"order_id": orderId,
		}).Infof("[skip-orderid] already paid %d", orderId)
		return nil
	}

	err = s.PaymentRepository.WithTransaction(ctx, func(tx *gorm.DB) error {
		err = s.PaymentRepository.UpdateStatus(ctx, tx, orderId, statusID)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"status_id": statusID,
				"order_id":  orderId,
				"err":       err.Error(),
			}).Errorf("s.PaymentRepository.UpdateStatus(ctx, orderId, statusID)")
			return err
		}

		// adding retry mechanism
		err = retryPublishEvent(maxRetryPublish, func() error {
			return s.Publisher.PublishPaymentSuccess(orderId)
		})

		if err != nil {
			failedEvent := models.FailedEvents{
				OrderID:    orderId,
				FailedType: constant.FailedPublishEventPaymentSuccess,
				Status:     constant.FailedPublishEventStatusNeedToCheck,
				Notes:      err.Error(),
				CreateTime: time.Now(),
			}

			// dead letter table
			publishErr := s.PaymentRepository.SaveFailedPublishEvent(ctx, failedEvent)
			if publishErr != nil {
				log.Logger.WithFields(logrus.Fields{
					"failed_param": failedEvent,
				}).WithError(publishErr)
			}

			log.Logger.WithFields(logrus.Fields{
				"order_id": orderId,
				"err":      err.Error(),
			}).Errorf("s.Publisher.PublishPaymentSuccess(orderId)")
			return err
		}

		return nil
	})

	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"order_id": orderId,
			"err":      err.Error(),
		}).Errorf("s.PaymentRepository.WithTransaction(ctx, func(tx *gorm.DB) error")
		return err
	}

	return nil
}

// CheckPaymentAmountByOrderID retrieves the payment amount for a given order ID.
// ctx: Context for managing request-scoped values.
// orderID: ID of the order.
func (s *paymentService) CheckPaymentAmountByOrderID(ctx context.Context, orderID int64) (float64, error) {
	amount, err := s.PaymentRepository.GetPaymentAmountByOrderID(ctx, orderID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"order_id": orderID,
		}).Errorf("error occurred on CheckPaymentAmountByOrderID(ctx context.Context, orderID int64): %s", err)
		return 0, err
	}

	return amount, nil
}

// SaveAnomaly saves a payment anomaly record to the database.
// ctx: Context for managing request-scoped values.
// param: Payment anomaly model containing anomaly details.
func (s *paymentService) SaveAnomaly(ctx context.Context, param models.PaymentAnomaly) error {
	err := s.PaymentRepository.SavePaymentAnomaly(ctx, param)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param": param,
		}).Errorf("error occured on PaymentService.SaveAnomaly(ctx context.Context, param models.PaymentAnomaly) %v", err)
		return err
	}
	return nil
}

// retryPublishEvent retries publishing an event up to a maximum number of attempts.
// max: Maximum number of retries.
// fn: Function to execute for publishing the event.
func retryPublishEvent(max int, fn func() error) error {
	var err error

	for i := range max {
		err = fn()
		if err == nil {
			return nil
		}

		wait := time.Duration(math.Pow(2, float64(1))) * time.Second
		log.Logger.Printf("retry: %d, Error: %s, retrying in %d secondss..", i+1, err, wait)
		time.Sleep(wait)
	}

	return err
}

func (s *paymentService) SavePaymentRequest(ctx context.Context, param models.PaymentRequests) error {
	err := s.PaymentRepository.SavePaymentRequest(ctx, param)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param": param,
		}).Errorf("error occurred on PaymentService.SavePaymentRequest(ctx context.Context, param models.PaymentRequests) %v", err)
	}
	return nil
}
