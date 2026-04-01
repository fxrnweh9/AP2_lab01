package main

import (
	"log"
	"payment-service/internal/app"
	"payment-service/internal/repository"
	"payment-service/internal/transport/http"
	"payment-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := app.NewDB("postgres://payment_user:1234@localhost:5432/payment_db")
	if err != nil {
		log.Fatal(err)
	}

	paymentRepo := repository.NewPaymentRepository(db)
	paymentUC := usecase.NewPaymentUseCase(paymentRepo)
	paymentHandler := http.NewPaymentHandler(paymentUC)

	r := gin.Default()
	paymentHandler.RegisterRoutes(r)

	log.Println("Payment Service running on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
