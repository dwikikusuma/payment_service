package routes

import (
	"github.com/gin-gonic/gin"
	"payment_service/cmd/handler"
	"payment_service/middleware"
)

func SetupRoutes(router *gin.Engine, paymentHandler handler.PaymentHandler, jwtSecret string) {
	router.Use(middleware.RequestLogger())
	authMiddleware := middleware.AuthMiddleware(jwtSecret)
	router.Use(authMiddleware)
}
