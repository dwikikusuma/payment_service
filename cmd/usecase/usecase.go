package usecase

import "order_service/cmd/service"

type OrderUseCase struct {
	OrderService service.OrderService
}

func NewOrderUseCase(orderService service.OrderService) *OrderUseCase {
	return &OrderUseCase{
		OrderService: orderService,
	}
}
