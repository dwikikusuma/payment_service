package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"payment_service/cmd/repository"
	"payment_service/infra/constant"
	"payment_service/infra/log"
)

type PaymentService interface {
	ProcessPaymentSuccess(ctx context.Context, orderId int64, status string) error
	CheckPaymentAmountByOrderID(ctx context.Context, orderID int64) (float64, error)
}

type paymentService struct {
	PaymentRepository repository.PaymentRepository
	Publisher         repository.PaymentEventPublisher
}

func NewPaymentService(paymentRepo repository.PaymentRepository, publisher repository.PaymentEventPublisher) PaymentService {
	return &paymentService{
		PaymentRepository: paymentRepo,
		Publisher:         publisher,
	}
}

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

		err = s.Publisher.PublishPaymentSuccess(orderId)
		if err != nil {
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
