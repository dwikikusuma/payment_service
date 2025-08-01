package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"payment_service/cmd/handler"
	"payment_service/cmd/repository"
	"payment_service/cmd/resource"
	"payment_service/cmd/service"
	"payment_service/cmd/usecase"
	"payment_service/config"
	"payment_service/infra/constant"
	"payment_service/infra/log"
	"payment_service/internalGrpc"
	internalKafka "payment_service/kafka"
	"payment_service/models"
	"payment_service/routes"
)

func main() {
	cfg := config.LoadConfig(
		config.WithConfigFolder([]string{"./files/config"}),
		config.WithConfigFile("config"),
		config.WithConfigType("yaml"),
	)
	fmt.Println("Configuration loaded successfully:", cfg)
	log.SetupLogger()

	db := resource.InitDB(&cfg)
	redis := resource.InitRedis(&cfg)
	writer := internalKafka.NewWriter(cfg.KafkaConfig.Broker, cfg.KafkaConfig.KafkaTopics[constant.KafkaTopicPaymentSuccess])

	constant.MapStatusFromDB(db)

	userClient := internalGrpc.NewUserClient()

	paymentRepo := repository.NewPaymentRepository(db, redis)
	paymentPublisher := repository.NewKafkaEventPublisher(writer)
	paymentService := service.NewPaymentService(paymentRepo, paymentPublisher)
	paymentUseCase := usecase.NewPaymentUseCase(paymentService)
	paymentHandler := handler.NewHandler(paymentUseCase, cfg.PGAConfig.WebhookToken)

	xenditRepository := repository.NewXenditClient(cfg.PGAConfig.ApiKey)
	xenditService := service.NewXenditService(paymentRepo, xenditRepository, *userClient)
	xenditUseCase := usecase.NewXenditUseCase(xenditService)

	schedulerService := service.SchedulerService{
		PaymentService:    paymentService,
		XenditClient:      xenditRepository,
		PublisherService:  paymentPublisher,
		PaymentRepository: paymentRepo,
		UserClient:        *userClient,
	}

	schedulerService.StartCheckPendingInvoice()
	schedulerService.StartProcessPaymentRequest()
	schedulerService.StartProcessFailedPaymentRequest()
	schedulerService.StartSweepingExpiredPendingPayments()

	internalKafka.StartKafkaConsumer(cfg.KafkaConfig.Broker, cfg.KafkaConfig.KafkaTopics[constant.KafkaTopicOrderCreated],
		func(event models.OrderCreatedEvent) {
			if cfg.Toggle.DisableCreatePaymentInvoiceDirectly {
				if err := paymentUseCase.ProcessPaymentRequest(context.Background(), event); err != nil {
					log.Logger.WithFields(logrus.Fields{
						"message": "failed to handle order created event",
						"err":     err.Error(),
					}).Error("paymentUseCase.ProcessPaymentRequest(context.Background(), event)")
				}
			} else {
				if err := xenditUseCase.CreateInvoice(context.Background(), event); err != nil {
					log.Logger.WithFields(logrus.Fields{
						"message": "failed to handle order created event",
						"err":     err.Error(),
					}).Error("xenditUseCase.CreateInvoice(context.Background(), event)")
				}
			}

		})
	router := gin.Default()
	routes.SetupRoutes(router, paymentHandler, "my")
	_ = router.Run(":" + cfg.App.Port)
	fmt.Println("Server running on port:", cfg.App.Port)
}
