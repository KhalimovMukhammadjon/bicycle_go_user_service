package grpc

import (
	"bicycle/user_service/config"
	"bicycle/user_service/genproto/user_service"
	"bicycle/user_service/grpc/client"
	"bicycle/user_service/grpc/service"
	"bicycle/user_service/pkg/logger"
	"bicycle/user_service/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(cfg config.Config, log logger.LoggerI, strg storage.StoragI, svcs client.ServiceManagerI) (grpcServer *grpc.Server) {
	grpcServer = grpc.NewServer()

	user_service.RegisterUserServiceServer(grpcServer, service.NewUserService(cfg, log, strg, svcs))

	reflection.Register(grpcServer)
	return
}
