package service

import (
	"context"
	"fmt"
	"payment_service/cmd/repository"
	"time"
)

type SchedulerService struct {
	PaymentRepository repository.PaymentRepository
	XenditClient      repository.XenditClient
	PublisherService  repository.PaymentEventPublisher
	PaymentService    PaymentService
}

func (s *SchedulerService) StartCheckPendingInvoice() {
	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for range ticker.C {
			ctx := context.Background()
			listPendingPayment, err := s.PaymentRepository.GetPendingPayment(ctx)
			if err != nil {
				fmt.Println("s.PaymentRepository.GetPendingPayment got error:", err)
				continue
			}

			for _, pendingInvoice := range listPendingPayment {
				invoiceStatus, err := s.XenditClient.CheckInvoiceStatus(ctx, pendingInvoice.ExternalID)
				if err != nil {
					fmt.Printf("s.XenditService.CheckInvoiceStatus got error: %s", err)
					continue
				}

				if invoiceStatus == "Paid" {
					err = s.PaymentService.ProcessPaymentSuccess(ctx, pendingInvoice.OrderID, "PAID")
					if err != nil {
						fmt.Printf("s.PaymentService.ProcessPaymentSuccess got error: %s", err)
						continue
					}
				}
			}
		}
	}()
}
