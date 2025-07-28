package service

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"payment_service/cmd/repository"
	"payment_service/infra/log"
	"payment_service/internalGrpc"
	"payment_service/models"
)

type XenditService interface {
	CreateInvoice(ctx context.Context, param models.OrderCreatedEvent) error
}

type xenditService struct {
	UserClient   internalGrpc.UserClient
	Database     repository.PaymentRepository
	XenditClient repository.XenditClient
}

func NewXenditService(db repository.PaymentRepository, xenditClient repository.XenditClient, userClient internalGrpc.UserClient) XenditService {
	return &xenditService{
		UserClient:   userClient,
		Database:     db,
		XenditClient: xenditClient,
	}
}

func (s *xenditService) CreateInvoice(ctx context.Context, param models.OrderCreatedEvent) error {
	externalID := fmt.Sprintf("order-%v", param.OrderID)
	userInfo, err := s.UserClient.GetUserByUserId(ctx, param.UserID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"message": "error getting user info",
			"exc":     "error occurred on s.UserClient.GetUserByUserId()",
			"err":     err.Error(),
			"userID":  param.UserID,
		}).Errorf("s.UserClient.GetUserByUserId(ctx, param.UserID)")
		return err
	}

	xenditReq := models.XenditInvoiceRequest{
		ExternalID:  externalID,
		Amount:      param.Amount,
		Description: fmt.Sprintf("Pembarayan Order %d", param.OrderID),
		PayerEmail:  userInfo.Email,
	}

	invoiceDetail, err := s.XenditClient.CrateInvoice(ctx, xenditReq)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"message": "error creating invoice to xendit",
			"exc":     "error occurred on xenditService.CreateInvoice()",
			"err":     err.Error(),
		})
		return err
	}

	paymentModel := models.Payment{
		OrderID:     param.OrderID,
		UserID:      param.UserID,
		ExternalID:  externalID,
		Amount:      param.Amount,
		Status:      "PENDING",
		ExpiredTime: invoiceDetail.ExpiryDate,
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
