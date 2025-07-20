package service

import (
	"context"
	"fmt"
	"log"
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
				log.Printf("Processing payment request for OrderID: %d, UserID: %d, Amount: %.2f", req.OrderID, req.UserID, req.Amount)
				externalID := fmt.Sprintf("order-%d", req.OrderID)
				paymentInfo, getPaymentErr := s.PaymentRepository.GetPaymentInfoByOrderID(ctx, req.OrderID)
				if getPaymentErr != nil {
					fmt.Printf("s.PaymentRepository.GetPaymentInfoByOrderID got error: %s", getPaymentErr)
					continue
				}

				if paymentInfo.ID != 0 {
					fmt.Printf("Payment already exists for OrderID: %d, skipping invoice creation", req.OrderID)
					updateErr := s.PaymentRepository.UpdateSuccessPaymentRequest(ctx, req.OrderID)
					if updateErr != nil {
						fmt.Printf("s.PaymentRepository.UpdateSuccessPaymentRequest got error: %s", updateErr)
						continue
					}
					continue
				}

				invoiceDetail, err := s.XenditClient.CrateInvoice(ctx, models.XenditInvoiceRequest{
					ExternalID:  externalID,
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

				updateErr := s.PaymentRepository.UpdateSuccessPaymentRequest(ctx, req.OrderID)
				if updateErr != nil {
					fmt.Printf("s.PaymentRepository.UpdateSuccessPaymentRequest got error: %s", updateErr)
					continue
				}

				savePaymentErr := s.PaymentRepository.SavePayment(ctx, models.Payment{
					OrderID:     req.OrderID,
					UserID:      req.UserID,
					Amount:      req.Amount,
					ExternalID:  externalID,
					Status:      "Pending",
					CreateTime:  time.Now(),
					ExpiredTime: invoiceDetail.ExpiryDate,
				})
				if savePaymentErr != nil {
					fmt.Printf("s.PaymentRepository.SavePayment got error: %s", savePaymentErr)
					continue
				}
			}
			time.Sleep(5 * time.Second) // Sleep to avoid busy loop
		}
	}(context.Background())
}

func (s *SchedulerService) StartProcessFailedPaymentRequest() {
	go func(ctx context.Context) {
		for {
			var failedPaymentReq []models.PaymentRequests
			err := s.PaymentRepository.GetFailedPaymentRequest(ctx, &failedPaymentReq)
			if err != nil {
				fmt.Printf("s.PaymentRepository.GetFailedPaymentRequest got error: %s", err)
				time.Sleep(5 * time.Second)
				continue
			}

			for _, req := range failedPaymentReq {
				updatePendingErr := s.PaymentRepository.UpdatePendingPaymentRequest(ctx, req.ID)
				if updatePendingErr != nil {
					fmt.Printf("s.PaymentRepository.UpdatePendingPaymentRequest got error: %s", updatePendingErr)
					updatedFailedErr := s.PaymentRepository.UpdateFailedPaymentRequest(ctx, req.OrderID, "Failed to update to pending")
					if updatedFailedErr != nil {
						fmt.Printf("s.PaymentRepository.UpdateFailedPaymentRequest got error: %s", updatedFailedErr)
					}
					continue
				}
			}
			time.Sleep(5 * time.Second) // Sleep to avoid busy loop
		}
	}(context.Background())
}

func (s *SchedulerService) StartSweepingExpiredPendingPayments() {
	go func(ctx context.Context) {
		for {
			fmt.Println("Checking for expired pending payments...")
			pendingPayments, err := s.PaymentRepository.GetExpiredPendingPayments(ctx)
			if err != nil {
				fmt.Printf("s.PaymentRepository.GetExpiredPendingPayments got error: %s", err)
				time.Sleep(5 * time.Second)
				continue
			}
			for _, pendingPayment := range pendingPayments {
				fmt.Printf("[updating] Processing expired payment for OrderID: %d\n", pendingPayment.OrderID)
				updatingErr := s.PaymentRepository.MarkExpiredPayments(ctx, pendingPayment.ID)
				if updatingErr != nil {
					fmt.Printf("s.PaymentRepository.MarkExpiredPayments got error: %s", updatingErr)
				} else {
					fmt.Printf("[updated] Successfully marked payment as expired for OrderID: %d\n", pendingPayment.OrderID)
				}
			}
			time.Sleep(10 * time.Second) // Sleep to avoid busy loop
		}
	}(context.Background())
}
