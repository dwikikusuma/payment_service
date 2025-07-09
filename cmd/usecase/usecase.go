package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"payment_service/cmd/service"
	"payment_service/infra/log"
	"payment_service/models"
	"strconv"
	"strings"
)

type PaymentUseCase interface {
	ProcessPaymentWebhook(ctx context.Context, payload models.XenditWebhookPayload) error
}

type paymentUseCase struct {
	PaymentService service.PaymentService
}

func NewPaymentUseCase(paymentService service.PaymentService) PaymentUseCase {
	return &paymentUseCase{
		PaymentService: paymentService,
	}
}

func (uc *paymentUseCase) ProcessPaymentWebhook(ctx context.Context, payload models.XenditWebhookPayload) error {
	switch payload.Status {
	case "PAID":
		orderID := extractOrderID(payload.ExternalID)
		// validate amount
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

		if amount != payload.Amount {
			errStr := fmt.Sprintf("webhook amount missmatch: expected %.2f, got %.2f", amount, payload.Amount)
			log.Logger.WithFields(logrus.Fields{
				"payload": payload,
			}).Error(errStr)
			return errors.New(errStr)
		}

		err = uc.PaymentService.ProcessPaymentSuccess(ctx, orderID, payload.Status)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"external_id": payload.ExternalID,
				"status":      payload.Status,
			}).Errorf("uc.PaymentService.ProcessPaymentSuccess(ctx, orderID, payload.Status)")
		}
	default:
		log.Logger.WithFields(logrus.Fields{
			"external_id": payload.ExternalID,
			"status":      payload.Status,
		}).Errorf("invalid status: %s", payload.Status)
		return errors.New("invalid status")
	}
	return nil
}

func extractOrderID(externalID string) int64 {
	orderIDStr := strings.TrimPrefix(externalID, "order-")
	orderID, err := strconv.ParseInt(orderIDStr, 64, 10)
	if err != nil {
		return 0
	}

	return orderID
}
