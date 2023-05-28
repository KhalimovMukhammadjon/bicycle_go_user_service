package storage

import (
	"bicycle/bicycle_go_user_service/genproto/user_service"
	"context"
)

type StoragI interface {
	CloseDB()
	User() UserRepoI
}

type UserRepoI interface {
	Create(ctx context.Context, req *user_service.CreateUserRequest) (resp *user_service.PrimaryKey, err error)
	GetById(ctx context.Context, req *user_service.PrimaryKey) (resp *user_service.User, err error)
	GetList(ctx context.Context, req *user_service.GetAllUserRequest) (resp *user_service.GetAllUserResponse, err error)
	// Update(ctx context.Context, req *user_service.PrimaryKey) error
	Delete(ctx context.Context, req *user_service.PrimaryKey) (err error)
	// GetUserByPhone(ctx context.Context, req *user_service.PhoneNumber) (resp *user_service.Checker, err error)
	GetUserByPhone(ctx context.Context, req *user_service.PrimaryKey) (resp *user_service.User, err error)
}
