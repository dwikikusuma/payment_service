package service

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"payment_service/cmd/repository"
	"payment_service/infra/log"
	"payment_service/models"
)

type XenditService interface {
	CreateInvoice(ctx context.Context, param models.OrderCreatedEvent) error
}

type xenditService struct {
	Database     repository.PaymentRepository
	XenditClient repository.XenditClient
}

func NewXenditService(db repository.PaymentRepository, xenditClient repository.XenditClient) XenditService {
	return &xenditService{
		Database:     db,
		XenditClient: xenditClient,
	}
}

func (s *xenditService) CreateInvoice(ctx context.Context, param models.OrderCreatedEvent) error {
	externalID := fmt.Sprintf("order-%v", param.OrderID)

	xenditReq := models.XenditInvoiceRequest{
		ExternalID:  externalID,
		Amount:      param.Amount,
		Description: fmt.Sprintf("Pembarayan Order %d", param.OrderID),
		PayerEmail:  "test@test.com",
	}

	_, err := s.XenditClient.CrateInvoice(ctx, xenditReq)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"message": "error creating invoice to xendit",
			"exc":     "error occurred on xenditService.CreateInvoice()",
			"err":     err.Error(),
		})
		return err
	}

	paymentModel := models.Payment{
		OrderID:    param.OrderID,
		UserID:     param.UserID,
		ExternalID: externalID,
		Amount:     param.Amount,
		Status:     "PENDING",
	}

	if err = s.Database.SavePayment(ctx, paymentModel); err != nil {
		log.Logger.WithFields(logrus.Fields{
			"message": "error saving payment",
			"exc":     "error occurred on s.Database.SavePayment(ctx, paymentModel)",
			"err":     err.Error(),
			"data":    paymentModel,
		}).Errorf("s.Database.SavePayment(ctx, paymentModel)")
		return err
	}

	return nil
}
