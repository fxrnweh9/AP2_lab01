package grpc

import (
	"order-service/internal/usecase"

	pb "github.com/fxrnweh9/proto-contracts/orderpb"
)

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
	uc *usecase.OrderUseCase
}

func NewOrderServer(uc *usecase.OrderUseCase) *OrderServer {
	return &OrderServer{uc: uc}
}
func (s *OrderServer) SubscribeToOrderUpdates(
	req *pb.OrderRequest,
	stream pb.OrderService_SubscribeToOrderUpdatesServer,
) error {

	orderID := req.OrderId

	for order := range s.uc.WatchOrder(orderID) {
		err := stream.Send(&pb.OrderStatusUpdate{
			OrderId: order.ID,
			Status:  string(order.Status),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
