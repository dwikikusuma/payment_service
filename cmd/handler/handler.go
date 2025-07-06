package handler

import "payment_service/cmd/usecase"

type PaymentHandler struct {
	PaymentUseCase usecase.PaymentUseCase
}

func NewHandler(paymentUseCase usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{
		PaymentUseCase: paymentUseCase,
	}
}
