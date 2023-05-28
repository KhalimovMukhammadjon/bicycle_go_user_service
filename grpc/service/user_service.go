package service

import (
	"bicycle/bicycle_go_user_service/config"
	"bicycle/bicycle_go_user_service/genproto/user_service"
	"bicycle/bicycle_go_user_service/grpc/client"
	"bicycle/bicycle_go_user_service/pkg/logger"
	"bicycle/bicycle_go_user_service/storage"
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
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
	o.log.Info("---Bicycle--->", logger.Any("req", req))

	pKey, err := o.strg.User().Create(context.Background(), req)
	if err != nil {
		o.log.Error("!!!CreateBook!!!", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	resp, err = o.strg.User().GetById(ctx, pKey)
	fmt.Println("err", err)

	return
}

// func (u userService) GetById(ctx context.Context, req *user_service.PrimaryKey) (resp *user_service.User, err error) {
// 	u.log.Info("---GetBicycle--->", logger.Any("req", req))
// 	fmt.Println(req.Id)

// 	resp, err = u.strg.User().GetById(ctx, req)

// 	if err != nil {
// 		u.log.Error("GetBicycle--->", logger.Error(err))
// 		return nil, status.Error(codes.InvalidArgument, err.Error())
// 	}

// 	return resp, err
// }

func (u userService) Delete(ctx context.Context, req *user_service.PrimaryKey) (resp *empty.Empty, err error) {
	u.log.Info("---Delete--->", logger.Any("req", req))

	resp = &empty.Empty{}

	err = u.strg.User().Delete(ctx, req)
	if err != nil {
		u.log.Error("!!!DeleteBook--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return resp, nil
}

// func (u userService) PhoneChecker(ctx context.Context, req *user_service.PhoneNumber) (resp *user_service.Checker, err error) {
// 	u.log.Info("---SMS code Checker--->", logger.Any("req", req))

// 	fmt.Println("Request--------------->>>>>>>>", req.PhoneNumber)

// 	user, err := u.strg.User().GetUserByPhone(ctx, req)
// 	if err != nil {
// 		u.log.Error("!!!PhoneChecker--->", logger.Error(err))
// 		return resp, status.Error(codes.InvalidArgument, err.Error())
// 	}

// 	if user != nil {
// 		resp.Check = true
// 	}
// 	fmt.Println("Test for user--------------->>>>>>>>", resp)

// 	return resp, nil
// }

func (u userService) GetById(ctx context.Context, req *user_service.PrimaryKey) (resp *user_service.User, err error) {
	u.log.Info("---GetBicycle--->", logger.Any("req", req))
	fmt.Println(req.Id)

	resp, err = u.strg.User().GetById(ctx, req)

	if err != nil {
		u.log.Error("GetBicycle--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return resp, err
}
