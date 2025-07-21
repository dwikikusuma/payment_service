package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"payment_service/cmd/service"
	"payment_service/infra/constant"
	"payment_service/infra/log"
	"payment_service/models"
	"payment_service/pdf"
	"strconv"
	"strings"
	"time"
)

// PaymentUseCase defines the interface for payment-related use case operations.
type PaymentUseCase interface {
	// ProcessPaymentWebhook processes a payment webhook payload.
	// ctx: Context for managing request-scoped values.
	// payload: XenditWebhookPayload containing webhook data.
	ProcessPaymentWebhook(ctx context.Context, payload models.XenditWebhookPayload) error

	// ProcessPaymentRequest processes a payment request.
	// ctx: Context for managing request-scoped values.
	// payload: OrderCreatedEvent containing details of the order.
	ProcessPaymentRequest(ctx context.Context, payload models.OrderCreatedEvent) error

	// DownloadPDFInvoice generates a PDF invoice for the given order ID.
	// ctx: Context for managing request-scoped values.
	// orderID: The ID of the order for which the invoice is generated.
	DownloadPDFInvoice(ctx context.Context, orderID int64) (string, error)
}

// paymentUseCase is the implementation of the PaymentUseCase interface.
type paymentUseCase struct {
	PaymentService service.PaymentService // Service for handling payment-related operations.
}

// NewPaymentUseCase creates a new instance of paymentUseCase.
// paymentService: Service for handling payment-related operations.
func NewPaymentUseCase(paymentService service.PaymentService) PaymentUseCase {
	return &paymentUseCase{
		PaymentService: paymentService,
	}
}

// ProcessPaymentWebhook processes a payment webhook payload.
// It validates the payment amount, handles anomalies, and updates the payment status.
// ctx: Context for managing request-scoped values.
// payload: XenditWebhookPayload containing webhook data.
func (uc *paymentUseCase) ProcessPaymentWebhook(ctx context.Context, payload models.XenditWebhookPayload) error {
	switch payload.Status {
	case "PAID":
		orderID := extractOrderID(payload.ExternalID)
		// Validate the payment amount.
		amount, err := uc.PaymentService.CheckPaymentAmountByOrderID(ctx, orderID)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"order_id":       orderID,
				"status":         payload.Status,
				"external_id":    payload.ExternalID,
				"webhook_amount": payload.Amount,
			})
			return err
		}

		// Check for amount mismatch and handle anomalies.
		if amount != payload.Amount {
			errStr := fmt.Sprintf("webhook amount missmatch: expected %.2f, got %.2f", amount, payload.Amount)
			paymentAnomaly := models.PaymentAnomaly{
				OrderID:     orderID,
				ExternalID:  payload.ExternalID,
				AnomalyType: constant.AnomalyTypeInvalidAmount,
				Notes:       errStr,
				Status:      constant.PaymentAnomalyStatusNeedToCheck,
				CreateTime:  time.Now(),
			}

			if err = uc.PaymentService.SaveAnomaly(ctx, paymentAnomaly); err != nil {
				log.Logger.WithFields(logrus.Fields{
					"payload":         payload,
					"payment_anomaly": paymentAnomaly,
				}).WithError(err)
			}
			log.Logger.WithFields(logrus.Fields{
				"payload": payload,
			}).Error(errStr)
			return errors.New(errStr)
		}

		// Process the payment success.
		err = uc.PaymentService.ProcessPaymentSuccess(ctx, orderID, "PAID")
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"external_id": payload.ExternalID,
				"status":      payload.Status,
			}).Errorf("uc.PaymentService.ProcessPaymentSuccess(ctx, orderID, payload.Status)")
		}
	default:
		// Handle invalid status.
		log.Logger.WithFields(logrus.Fields{
			"external_id": payload.ExternalID,
			"status":      payload.Status,
		}).Errorf("invalid status: %s", payload.Status)
		return errors.New("invalid status")
	}
	return nil
}

func (uc *paymentUseCase) ProcessPaymentRequest(ctx context.Context, payload models.OrderCreatedEvent) error {
	err := uc.PaymentService.SavePaymentRequest(ctx, models.PaymentRequests{
		OrderID:    payload.OrderID,
		Amount:     payload.Amount,
		UserID:     payload.UserID,
		Status:     "PENDING",
		CreateTime: time.Now(),
	})

	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"order_id": payload.OrderID,
			"amount":   payload.Amount,
			"user_id":  payload.UserID,
		}).Errorf("error occurred on ProcessPaymentRequest: %v", err)
		return err
	}

	return nil
}

// extractOrderID extracts the order ID from the external ID.
// externalID: The external ID string.
// Returns the extracted order ID as an int64.
func extractOrderID(externalID string) int64 {
	orderIDStr := strings.TrimPrefix(externalID, "order-")
	orderID, err := strconv.ParseInt(orderIDStr, 64, 10)
	if err != nil {
		return 0
	}

	return orderID
}

// DownloadPDFInvoice generates a PDF invoice for the given order ID.
// ctx: Context for managing request-scoped values.
// orderID: The ID of the order for which the invoice is generated.
// Returns the file path of the generated PDF invoice or an error if the operation fails.
func (uc *paymentUseCase) DownloadPDFInvoice(ctx context.Context, orderID int64) (string, error) {
	paymentDetail, err := uc.PaymentService.GetPaymentInfoByOrderID(ctx, orderID)
	if err != nil {
		return "", err
	}

	filePath := fmt.Sprintf("/fcproject/invoice_%d", orderID)
	err = pdf.GenerateInvoicePDF(paymentDetail, filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
