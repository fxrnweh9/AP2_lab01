package main

import (
	"context"
	"log"
	"net"
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
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	grpcLib "google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()

	// 1. db
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:pass@order-db:5432/orders_db?sslmode=disable"
	}
	db, err := app.NewDB(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Redis
	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis not reachable: %v", err)
	}

	// 3. Repo and client
	orderRepo := repository.NewOrderRepository(db)

	paymentAddr := os.Getenv("PAYMENT_GRPC_ADDR")
	if paymentAddr == "" {
		paymentAddr = "payment-service:50051"
	}
	paymentClient, err := app.NewPaymentGRPCClient(paymentAddr)
	if err != nil {
		log.Fatal(err)
	}

	// 4. UseCase
	orderUC := usecase.NewOrderUseCase(orderRepo, paymentClient, rdb)

	// 5. HTTP + Rate Limiter
	orderHandler := orderHandlerHttp.NewOrderHandler(orderUC)
	r := gin.Default()
	r.Use(orderHandlerHttp.RateLimiter(rdb, 10, time.Minute)) // ← используем rdb, не redisClient
	orderHandler.RegisterRoutes(r)

	// 6. gRPC
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

	// 7. Graceful Shutdown
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
