package routes

import (
	"github.com/gin-gonic/gin"
	"payment_service/cmd/handler"
	"payment_service/middleware"
)

func SetupRoutes(router *gin.Engine, paymentHandler handler.PaymentHandler, jwtSecret string) {
	router.Use(middleware.RequestLogger())
	router.POST("/v1/payment/webhook", paymentHandler.HandleXenditWebhook)
	router.GET("/v1/payment/invoice/:order_id/pdf", paymentHandler.HandleDownloadPDFInvoice)
}
