package repository

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"payment_service/infra/log"
	"payment_service/models"
)

func (r *paymentRepository) UpdateStatus(ctx context.Context, tx *gorm.DB, orderId int64, status int64) error {
	return tx.
		WithContext(ctx).
		Model(models.Payment{}).
		Where("order_id = ?", orderId).
		Update("status", status).
		Error
}

func (r *paymentRepository) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := r.Database.Begin().WithContext(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *paymentRepository) SavePayment(ctx context.Context, model models.Payment) error {
	return r.Database.Table("payments").WithContext(ctx).Create(&model).Error
}

func (r *paymentRepository) IsAlreadyPaid(ctx context.Context, orderID int64) (bool, error) {
	var result models.Payment
	err := r.Database.Table("payments").WithContext(ctx).Where("order_id = ?", orderID).First(&result).Error
	if err != nil {
		return false, err
	}
	return result.Status == "PAID", nil
}

func (r *paymentRepository) GetPaymentAmountByOrderID(ctx context.Context, orderID int64) (float64, error) {
	var result models.Payment
	err := r.Database.Table("payments").WithContext(ctx).Where("order_id = ?", orderID).First(&result).Error
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"order_id": orderID,
		}).Errorf("error occurred on GetPaymentAmountByOrderID(ctx context.Context, orderID int64): %s", err)
		return 0, err
	}

	return result.Amount, nil
}

func (r *paymentRepository) SavePaymentAnomaly(ctx context.Context, param models.PaymentAnomaly) error {
	err := r.Database.Table("payment_anomalies").WithContext(ctx).Create(param).Error
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param": param,
		}).Errorf("error occurred on SavePaymentAnomaly(ctx context.Context, param models.PaymentAnomaly) %v", err)
	}
	return nil
}
