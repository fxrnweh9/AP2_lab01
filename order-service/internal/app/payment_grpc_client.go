package app

import (
	"context"
	"time"

	"order-service/internal/domain"

	pb "github.com/fxrnweh9/proto-contracts/paymentpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentGRPCClient struct {
	client pb.PaymentServiceClient
}

func NewPaymentGRPCClient(addr string) (*PaymentGRPCClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	if err != nil {
		return nil, err
	}

	return &PaymentGRPCClient{
		client: pb.NewPaymentServiceClient(conn),
	}, nil
}

func (p *PaymentGRPCClient) ProcessPayment(orderID string, amount int64) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := p.client.ProcessPayment(ctx, &pb.PaymentRequest{
		OrderId: orderID,
		Amount:  amount,
	})

	if err != nil {
		// Конвертируем gRPC ошибку в domain ошибку
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.FailedPrecondition:
				return "", domain.ErrPaymentLimitExceeded
			case codes.Unavailable, codes.DeadlineExceeded:
				return "", domain.ErrServiceUnavailable
			}
		}
		return "", err
	}

	return resp.Status, nil
}
