package postgres

import (
	"bicycle/bicycle_go_user_service/genproto/user_service"
	"bicycle/bicycle_go_user_service/storage"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) storage.UserRepoI {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) Create(ctx context.Context, req *user_service.CreateUserRequest) (resp *user_service.PrimaryKey, err error) {
	query :=
		`
		INSERT INTO users
		(id,first_name,last_name,phone_number) 
		VALUES(
			$1,
			$2,
			$3,
			$4
		)
	`
	uuid, err := uuid.NewRandom()
	if err != nil {
		return resp, err
	}

	_, err = u.db.Exec(ctx, query,
		uuid.String(),
		req.FirstName,
		req.LastName,
		req.PhoneNumber,
	)
	if err != nil {
		return resp, err
	}

	resp = &user_service.PrimaryKey{
		Id: uuid.String(),
	}

	return resp, nil
}

func (u *userRepo) GetById(ctx context.Context, req *user_service.PrimaryKey) (resp *user_service.User, err error) {
	resp = &user_service.User{}

	query := `
		SELECT 
			id,
			first_name,
			last_name,
			phone_number
		FROM users
		WHERE id = $1
	`
	err = u.db.QueryRow(ctx, query, req.Id).Scan(
		&resp.Id,
		&resp.FirstName,
		&resp.LastName,
		&resp.PhoneNumber,
	)

	if err != nil {
		return resp, err
	}

	return resp,nil
}

func (u *userRepo) GetList(ctx context.Context, req *user_service.GetAllUserRequest) (resp *user_service.GetAllUserResponse, err error) {
	// var (
	// 	params (map[string]interface{})

	// )

	return resp, nil
}

func (u *userRepo) Delete(ctx context.Context, req *user_service.PrimaryKey) (err error) {
	query := `DELETE FROM users WHERE id = $1`

	_, err = u.db.Exec(ctx, query, req.Id)
	if err != nil {
		return err
	}

	return nil
}
