package repository

import (
	"context"
	"gorm.io/gorm"
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
