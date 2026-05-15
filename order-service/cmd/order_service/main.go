package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"order-service/internal/app"
	"order-service/internal/repository"
	orderGrpc "order-service/internal/transport/grpc"
	orderHandlerHttp "order-service/internal/transport/http"
	"order-service/internal/usecase"

	pb "github.com/fxrnweh9/proto-contracts/orderpb"

	"github.com/gin-gonic/gin"
	grpcLib "google.golang.org/grpc"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:pass@order-db:5432/orders_db?sslmode=disable"
	}

	db, err := app.NewDB(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	orderRepo := repository.NewOrderRepository(db)

	paymentAddr := os.Getenv("PAYMENT_GRPC_ADDR")
	if paymentAddr == "" {
		paymentAddr = "payment-service:50051"
	}

	paymentClient, err := app.NewPaymentGRPCClient(paymentAddr)
	if err != nil {
		log.Fatal(err)
	}

	orderUC := usecase.NewOrderUseCase(orderRepo, paymentClient)

	orderHandler := orderHandlerHttp.NewOrderHandler(orderUC)
	r := gin.Default()
	orderHandler.RegisterRoutes(r)

	orderServer := orderGrpc.NewOrderServer(orderUC)

	go func() {
		lis, err := net.Listen("tcp", ":50052")
		if err != nil {
			log.Fatal(err)
		}
		grpcServer := grpcLib.NewServer()
		pb.RegisterOrderServiceServer(grpcServer, orderServer)
		log.Println("Order gRPC Server running on :50052")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Order Service REST running on :8080")
		if err := r.Run(":8080"); err != nil {
			log.Fatal(err)
		}
	}()

	<-quit
	log.Println("Shutting down order-service...")
}
