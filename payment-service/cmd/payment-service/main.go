package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"payment-service/internal/app"
	"payment-service/internal/repository"
	grpcHandler "payment-service/internal/transport/grpc"
	"payment-service/internal/usecase"

	pb "github.com/fxrnweh9/proto-contracts/paymentpb"

	"google.golang.org/grpc"

	grpcLib "google.golang.org/grpc"

	"github.com/joho/godotenv"
)

func LoggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {

	start := time.Now()

	res, err := handler(ctx, req)

	log.Printf("method=%s duration=%s error=%v",
		info.FullMethod,
		time.Since(start),
		err,
	)

	return res, err
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := app.NewDB(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewPaymentRepository(db)
	uc := usecase.NewPaymentUseCase(repo)
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
