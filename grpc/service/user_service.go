package service

import (
	"bicycle/user_service/config"
	"bicycle/user_service/genproto/user_service"
	"bicycle/user_service/grpc/client"
	"bicycle/user_service/pkg/logger"
	"bicycle/user_service/storage"
)

type userService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StoragI
	services client.ServiceManagerI
	user_service.UnimplementedUserServiceServer
}

func NewUserService(cfg config.Config, log logger.LoggerI, strg storage.StoragI, svcs client.ServiceManagerI) *userService {
	return &userService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: svcs,
	}
}
