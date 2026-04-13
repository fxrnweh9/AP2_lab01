package grpc

import (
	"context"

	pb "github.com/fxrnweh9/proto-contracts/paymentpb"
	"payment-service/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentServer struct {
	pb.UnimplementedPaymentServiceServer
	uc *usecase.PaymentUseCase
}

func NewPaymentServer(uc *usecase.PaymentUseCase) *PaymentServer {
	return &PaymentServer{uc: uc}
}

func (s *PaymentServer) ProcessPayment(
	ctx context.Context,
	req *pb.PaymentRequest,
) (*pb.PaymentResponse, error) {

	payment, err := s.uc.ProcessPayment(req.OrderId, req.Amount)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &pb.PaymentResponse{
		Status:        string(payment.Status),
		TransactionId: payment.TransactionID,
	}, nil
}
