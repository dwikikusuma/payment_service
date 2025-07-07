package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"payment_service/cmd/usecase"
	"payment_service/infra/log"
	"payment_service/models"
)

type PaymentHandler interface {
	HandleXenditWebhook(c *gin.Context)
}

type paymentHandler struct {
	PaymentUseCase usecase.PaymentUseCase
}

func NewHandler(paymentUseCase usecase.PaymentUseCase) PaymentHandler {
	return &paymentHandler{
		PaymentUseCase: paymentUseCase,
	}
}

func (h *paymentHandler) HandleXenditWebhook(c *gin.Context) {
	var model models.XenditWebhookPayload
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid payload",
		})
		return
	}

	if err := h.PaymentUseCase.ProcessPaymentWebhook(c.Request.Context(), model); err != nil {
		log.Logger.WithFields(logrus.Fields{
			"message": "failed to process payment",
			"err":     err.Error(),
		}).Error("h.PaymentUseCase.ProcessPaymentWebhook(c.Request.Context(), model)")

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "invalid payload",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})

}
