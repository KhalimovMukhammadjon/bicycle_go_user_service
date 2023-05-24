package client

import (
	"bicycle/bicycle_go_user_service/config"
	"bicycle/bicycle_go_user_service/genproto/order_service"

	"google.golang.org/grpc"
)

type ServiceManagerI interface {
	OrderService() order_service.OrderServiceClient
	ProductService() order_service.ProductServiceClient
}

type grpcClients struct {
	orderService   order_service.OrderServiceClient
	productService order_service.ProductServiceClient
}

func NewGrpcClient(cfg config.Config) (ServiceManagerI, error) {
	connUserService, err := grpc.Dial(
		cfg.OrderServiceHost+cfg.OrderServicePort,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return &grpcClients{
		orderService: order_service.NewOrderServiceClient(connUserService),
	}, nil
}

func (g *grpcClients) OrderService() order_service.OrderServiceClient {
	return g.orderService
}

func (g *grpcClients) ProductService() order_service.ProductServiceClient {
	return g.productService
}
