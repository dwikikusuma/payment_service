package usecase

import (
	"context"
	"github.com/sirupsen/logrus"
	"payment_service/cmd/service"
	"payment_service/infra/log"
	"payment_service/models"
)

type XenditUseCase interface {
	CreateInvoice(ctx context.Context, param models.OrderCreatedEvent) error
}

type xenditUseCase struct {
	svc service.XenditService
}

func NewXenditUseCase(svc service.XenditService) XenditUseCase {
	return &xenditUseCase{
		svc: svc,
	}
}

func (uc *xenditUseCase) CreateInvoice(ctx context.Context, param models.OrderCreatedEvent) error {
	if err := uc.svc.CreateInvoice(ctx, param); err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param": param,
		}).Errorf("error occurred on xenditUseCase.CreateInvoice(), %v ", err)
		return err
	}
	return nil
}
