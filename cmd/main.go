package main

import (
	"bicycle/bicycle_go_user_service/config"
	"bicycle/bicycle_go_user_service/grpc"
	"bicycle/bicycle_go_user_service/grpc/client"
	"bicycle/bicycle_go_user_service/pkg/logger"
	"bicycle/bicycle_go_user_service/storage/postgres"
	"context"
	"net"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	// ----------------------------------------------
	var loggerLevel = new(string)
	*loggerLevel = logger.LevelDebug
	switch cfg.Environment {
	case config.DebugMode:
		*loggerLevel = logger.LevelDebug
		gin.SetMode(gin.DebugMode)
	case config.TestMode:
		*loggerLevel = logger.LevelDebug
		gin.SetMode(gin.TestMode)
	default:
		*loggerLevel = logger.LevelInfo
		gin.SetMode(gin.ReleaseMode)

	}

	log := logger.NewLogger("app", *loggerLevel)
	defer func() {
		err := logger.Cleanup(log)
		if err != nil {
			return
		}
	}()

	// ----------------------------------------------

	store, err := postgres.NewPostgres(ctx, cfg)
	if err != nil {
		log.Panic("Error connect to postgresql: ", logger.Error(err))
		return
	}
	defer store.CloseDB()

	svcs, err := client.NewGrpcClient(cfg)
	if err != nil {
		log.Panic("client.NewGrpcClients", logger.Error(err))
	}

	grpcServer := grpc.SetUpServer(cfg, log, store, svcs)

	lis, err := net.Listen("tcp", cfg.ServicePort)
	if err != nil {
		log.Panic("net.Listen", logger.Error(err))
	}

	log.Info("GRPC: Server being started...", logger.String("port", cfg.ServicePort))

	if err := grpcServer.Serve(lis); err != nil {
		log.Panic("grpcServer.Serve", logger.Error(err))
	}
}
