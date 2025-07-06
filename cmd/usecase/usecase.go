package usecase

import "payment_service/cmd/service"

type PaymentUseCase struct {
	PaymentService service.PaymentService
}

func NewPaymentUseCase(paymentService service.PaymentService) *PaymentUseCase {
	return &PaymentUseCase{
		PaymentService: paymentService,
	}
}
