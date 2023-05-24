package service

import (
	"bicycle/bicycle_go_user_service/config"
	"bicycle/bicycle_go_user_service/genproto/user_service"
	"bicycle/bicycle_go_user_service/grpc/client"
	"bicycle/bicycle_go_user_service/pkg/logger"
	"bicycle/bicycle_go_user_service/storage"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (o userService) Create(ctx context.Context, req *user_service.CreateUserRequest) (resp *user_service.User, err error) {
	pKey, err := o.strg.User().Create(context.Background(), req)
	if err != nil {
		o.log.Error("!!!CreateBook!!!", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	resp, err = o.strg.User().GetById(ctx, pKey)

	return
}
