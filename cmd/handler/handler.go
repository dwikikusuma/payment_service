package handler

import "order_service/cmd/usecase"

type OrderHandler struct {
	OrderUseCase usecase.OrderUseCase
}

func NewHandler(orderUseCase usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		OrderUseCase: orderUseCase,
	}
}
