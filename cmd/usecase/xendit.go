package usecase

import (
	"context"
	"github.com/sirupsen/logrus"
	"payment_service/cmd/service"
	"payment_service/infra/log"
	"payment_service/models"
)

// XenditUseCase defines the interface for Xendit-related use case operations.
type XenditUseCase interface {
	// CreateInvoice creates an invoice based on the provided order event.
	// ctx: Context for managing request-scoped values.
	// param: OrderCreatedEvent containing details of the order.
	CreateInvoice(ctx context.Context, param models.OrderCreatedEvent) error
}

// xenditUseCase is the implementation of the XenditUseCase interface.
type xenditUseCase struct {
	svc service.XenditService // Service for interacting with Xendit operations.
}

// NewXenditUseCase creates a new instance of xenditUseCase.
// svc: XenditService instance for handling Xendit-related operations.
func NewXenditUseCase(svc service.XenditService) XenditUseCase {
	return &xenditUseCase{
		svc: svc,
	}
}

// CreateInvoice creates an invoice using the Xendit service.
// It logs any errors that occur during the process.
// ctx: Context for managing request-scoped values.
// param: OrderCreatedEvent containing details of the order.
func (uc *xenditUseCase) CreateInvoice(ctx context.Context, param models.OrderCreatedEvent) error {
	if err := uc.svc.CreateInvoice(ctx, param); err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param": param,
		}).Errorf("error occurred on xenditUseCase.CreateInvoice(), %v ", err)
		return err
	}
	return nil
}
