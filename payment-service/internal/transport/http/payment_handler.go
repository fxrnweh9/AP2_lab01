package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"payment-service/internal/domain"
	"payment-service/internal/usecase"
)

type PaymentHandler struct {
	uc *usecase.PaymentUseCase
}

func NewPaymentHandler(uc *usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{uc: uc}
}

type PaymentRequest struct {
	OrderID string `json:"order_id"`
	Amount  int64  `json:"amount"`
}

func (h *PaymentHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/payments", h.ProcessPayment)
	r.GET("/payments/:order_id", h.GetPayment)
}

func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.uc.ProcessPayment(req.OrderID, req.Amount)
	if err != nil && err != domain.ErrPaymentLimitExceeded {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         payment.Status,
		"transaction_id": payment.TransactionID,
	})
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	orderID := c.Param("order_id")
	payment, err := h.uc.GetPayment(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}

	c.JSON(http.StatusOK, payment)
}
