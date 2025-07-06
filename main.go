package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"payment_service/cmd/handler"
	"payment_service/cmd/repository"
	"payment_service/cmd/resource"
	"payment_service/cmd/service"
	"payment_service/cmd/usecase"
	"payment_service/config"
	"payment_service/infra/log"
	"payment_service/routes"
)

func main() {
	cfg := config.LoadConfig(
		config.WithConfigFolder([]string{"./files/config"}),
		config.WithConfigFile("config"),
		config.WithConfigType("yaml"),
	)
	log.SetupLogger()

	db := resource.InitDB(&cfg)
	redis := resource.InitRedis(&cfg)

	paymentRepo := repository.NewPaymentRepository(db, redis)
	paymentService := service.NewPaymentService(*paymentRepo)
	paymentUseCase := usecase.NewPaymentUseCase(*paymentService)
	paymentHandler := handler.NewHandler(*paymentUseCase)

	fmt.Println("Configuration loaded successfully:", cfg)

	router := gin.Default()
	routes.SetupRoutes(router, *paymentHandler, "my")
	_ = router.Run(":" + cfg.App.Port)
	fmt.Println("Server running on port:", cfg.App.Port)
}
