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

func (r *paymentRepository) SaveFailedPublishEvent(ctx context.Context, param models.FailedEvents) error {
	err := r.Database.Table("failed_events").WithContext(ctx).Create(param).Error
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param": param,
		}).WithError(err)
		return err
	}
	return nil
}

func (r *paymentRepository) GetPendingPayment(ctx context.Context) ([]models.Payment, error) {
	var result []models.Payment

	err := r.Database.Table("payment").WithContext(ctx).Where("status = ?", "Pending").Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *paymentRepository) SavePaymentRequest(ctx context.Context, param models.PaymentRequests) error {
	err := r.Database.Table("payment_requests").WithContext(ctx).Create(models.PaymentRequests{
		OrderID:    param.OrderID,
		UserID:     param.UserID,
		Amount:     param.Amount,
		Status:     param.Status,
		CreateTime: param.CreateTime,
	}).Error

	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param": param,
		}).Errorf("error occurred on SavePaymentRequest(ctx context.Context, param models.PaymentRequests) %v", err)

		return err
	}

	return nil
}

func (r *paymentRepository) GetPendingPaymentRequest(ctx context.Context, paymentRequests *[]models.PaymentRequests) error {
	err := r.Database.Table("payment_requests").WithContext(ctx).Where("status = ?", "Pending").Order("create_time ASC").Find(paymentRequests).Error
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Errorf("error occurred on GetPendingPaymentRequest(ctx context.Context, paymentRequests *[]models.PaymentRequests) %v", err)
		return err
	}
	return nil
}

func (r *paymentRepository) UpdateFailedPaymentRequest(ctx context.Context, orderID int64, notes string) error {
	err := r.Database.Table("payment_requests").WithContext(ctx).Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"status":      "Failed",
			"notes":       notes,
			"retry_count": gorm.Expr("retry_count + 1"),
			"update_time": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error

	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"order_id": orderID,
			"notes":    notes,
		}).Errorf("error occurred on UpdateFailedPaymentRequest(ctx context.Context, orderID int64, notes string) %v", err)
		return err
	}
	return nil
}

func (r *paymentRepository) UpdateSuccessPaymentRequest(ctx context.Context, orderID int64) error {
	err := r.Database.Table("payment_requests").WithContext(ctx).Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"status":      "Success",
			"update_time": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error

	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"order_id": orderID,
		}).Errorf("error occurred on UpdateSuccessPaymentRequest(ctx context.Context, orderID int64) %v", err)
		return err
	}
	return nil
}

func (r *paymentRepository) GetFailedPaymentRequest(ctx context.Context, requests *[]models.PaymentRequests) error {
	err := r.Database.Table("payment_requests").WithContext(ctx).
		Where("status = ?", "Failed").
		Where("retry_count < ?", 3).
		Order("create_time ASC").
		Find(requests).Error

	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Errorf("error occurred on GetFailedPaymentRequest(ctx context.Context, requests *[]models.PaymentRequests) %v", err)
		return err
	}

	return nil
}

func (r *paymentRepository) UpdatePendingPaymentRequest(ctx context.Context, paymentRequestID int64) error {
	err := r.Database.Table("payment_requests").WithContext(ctx).Where("id = ?", paymentRequestID).
		Updates(map[string]interface{}{
			"status":      "Pending",
			"update_time": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error

	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"payment_request_id": paymentRequestID,
		}).Errorf("error occurred on UpdatePendingPaymentRequest(ctx context.Context, paymentRequestID int64) %v", err)
		return err
	}

	return nil

}

func (r *paymentRepository) GetPaymentInfoByOrderID(ctx context.Context, orderID int64) (models.Payment, error) {
	var result models.Payment
	err := r.Database.Table("payments").WithContext(ctx).Where("order_id = ?", orderID).First(&result).Error
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"order_id": orderID,
		}).Errorf("error occurred on GetPaymentInfoByOrderID(ctx context.Context, orderID int64): %s", err)
		return models.Payment{}, err
	}
	return result, nil
}

func (r *paymentRepository) GetExpiredPendingPayments(ctx context.Context) ([]models.Payment, error) {
	var result []models.Payment
	err := r.Database.Table("payments").WithContext(ctx).
		Where("status = ?", "Pending").
		Where("expired_time < ?", gorm.Expr("CURRENT_TIMESTAMP")).
		Find(&result).Error

	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Errorf("error occurred on GetExpiredPendingPayments(ctx context.Context): %s", err)
		return nil, err
	}
	return result, nil
}

func (r *paymentRepository) MarkExpiredPayments(ctx context.Context, paymentId int64) error {
	err := r.Database.WithContext(ctx).
		Table("payments").
		Where("id = ?", paymentId).
		Updates(map[string]interface{}{
			"status":      "Expired",
			"update_time": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error

	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"payment_id": paymentId,
		}).Errorf("error occurred on MarkExpiredPayments(ctx context.Context, paymentId int64): %s", err)
		return err
	}
	return nil
}
