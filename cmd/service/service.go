package service

import "order_service/cmd/repository"

type OrderService struct {
	OrderRepository repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) *OrderService {
	return &OrderService{
		OrderRepository: orderRepo,
	}
}
