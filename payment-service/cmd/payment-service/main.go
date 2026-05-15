package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"payment-service/internal/app"
	"payment-service/internal/infrastructure/rabbitmq"
	"payment-service/internal/repository"
	grpcHandler "payment-service/internal/transport/grpc"
	"payment-service/internal/usecase"

	pb "github.com/fxrnweh9/proto-contracts/paymentpb"

	grpcLib "google.golang.org/grpc"

	"github.com/joho/godotenv"
)

func LoggingInterceptor(
	ctx context.Context,
	req any,
	info *grpcLib.UnaryServerInfo,
	handler grpcLib.UnaryHandler,
) (any, error) {
	start := time.Now()
	res, err := handler(ctx, req)
	log.Printf("method=%s duration=%s error=%v", info.FullMethod, time.Since(start), err)
	return res, err
}

func main() {
	_ = godotenv.Load()

	db, err := app.NewDB(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	// RabbitMQ producer
	producer, err := rabbitmq.NewProducer(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer producer.Close()

	publisher := rabbitmq.NewPublisherAdapter(producer)

	repo := repository.NewPaymentRepository(db)
	uc := usecase.NewPaymentUseCase(repo, publisher)
	server := grpcHandler.NewPaymentServer(uc)

	lis, err := net.Listen("tcp", ":"+os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpcLib.NewServer(
		grpcLib.UnaryInterceptor(LoggingInterceptor),
	)

	pb.RegisterPaymentServiceServer(grpcServer, server)
	log.Println("gRPC Payment Service running on", os.Getenv("GRPC_PORT"))

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
