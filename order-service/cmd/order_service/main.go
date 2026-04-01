package main

import (
	"log"
	"order-service/internal/app"
	"order-service/internal/repository"
	"order-service/internal/transport/http"
	"order-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := app.NewDB("postgres://order_user:1234@localhost:5432/order_db")
	if err != nil {
		log.Fatal(err)
	}

	orderRepo := repository.NewOrderRepository(db)
	paymentClient := app.NewPaymentClient("http://localhost:8081") // Payment Service URL

	orderUC := usecase.NewOrderUseCase(orderRepo, paymentClient)
	orderHandler := http.NewOrderHandler(orderUC)

	r := gin.Default()
	orderHandler.RegisterRoutes(r)

	log.Println("Order Service running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
