package service

import (
	"context"
	"fmt"
	"payment_service/cmd/repository"
	"payment_service/models"
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

func (s *SchedulerService) StartProcessPaymentRequest() {
	go func(ctx context.Context) {
		for {
			var paymentReq []models.PaymentRequests
			err := s.PaymentRepository.GetPendingPaymentRequest(ctx, &paymentReq)
			if err != nil {
				fmt.Printf("s.PaymentRepository.GetPendingPayment got error: %s", err)
				time.Sleep(5 * time.Second)
				continue
			}

			for _, req := range paymentReq {
				_, err = s.XenditClient.CrateInvoice(ctx, models.XenditInvoiceRequest{
					ExternalID:  fmt.Sprintf("order-%d", req.OrderID),
					Amount:      req.Amount,
					Description: fmt.Sprintf("Payment for order %d", req.OrderID),
					PayerEmail:  fmt.Sprintf("user%d@test.com", req.UserID),
				})

				if err != nil {
					errSavedFailed := s.PaymentRepository.UpdateFailedPaymentRequest(ctx, req.OrderID, err.Error())
					if errSavedFailed != nil {
						fmt.Printf("s.PaymentRepository.UpdateFailedPaymentRequest got error: %s", errSavedFailed)
					}
					continue
				}
			}
		}
	}(context.Background())
}
