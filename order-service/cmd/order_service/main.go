package main

import (
	"log"
	"net"

	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"order-service/internal/app"
	"order-service/internal/repository"
	orderGrpc "order-service/internal/transport/grpc"
	orderHandlerHttp "order-service/internal/transport/http"
	"order-service/internal/usecase"

	pb "github.com/fxrnweh9/proto-contracts/orderpb"

	"github.com/gin-gonic/gin"
	grpcLib "google.golang.org/grpc"
)

func main() {
	db, err := app.NewDB("postgres://order_user:1234@localhost:5432/order_db")
	if err != nil {
		log.Fatal(err)
	}

	orderRepo := repository.NewOrderRepository(db)

	paymentClient, err := app.NewPaymentGRPCClient("localhost:50051")
	if err != nil {
		log.Fatal(err)
	}

	orderUC := usecase.NewOrderUseCase(orderRepo, paymentClient)

	// REST
	orderHandler := orderHandlerHttp.NewOrderHandler(orderUC)
	r := gin.Default()
	orderHandler.RegisterRoutes(r)

	// gRPC server
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

	// REST runs last (blocking)
	log.Println("Order Service running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = ctx

}
