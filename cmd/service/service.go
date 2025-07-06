package service

import "payment_service/cmd/repository"

type PaymentService struct {
	PaymentRepository repository.PaymentRepository
}

func NewPaymentService(paymentRepo repository.PaymentRepository) *PaymentService {
	return &PaymentService{
		PaymentRepository: paymentRepo,
	}
}
