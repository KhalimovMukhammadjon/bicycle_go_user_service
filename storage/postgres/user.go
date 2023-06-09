package postgres

import (
	"bicycle/bicycle_go_user_service/genproto/user_service"
	"bicycle/bicycle_go_user_service/pkg/helper"
	"bicycle/bicycle_go_user_service/storage"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	return resp, nil
}

func (u *userRepo) GetList(ctx context.Context, req *user_service.GetAllUserRequest) (resp *user_service.GetAllUserResponse, err error) {
	resp = &user_service.GetAllUserResponse{}

	var (
		filter = " WHERE TRUE "
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query := `
		SELECT
			id,
			first_name,
			last_name,
			phone_number
		FROM users
	`

	if len(req.Search) > 0 {
		filter += " AND first_name ILIKE '%' || '" + req.Search + "' || '%' "
	}

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	query += filter + offset + limit

	rows, err := u.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user user_service.User
		err = rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.PhoneNumber,
		)
		if err != nil {
			return nil, err
		}
		resp.User = append(resp.User, &user)
	}
	return resp, nil
}

func (u *userRepo) Update(ctx context.Context, req *user_service.UpdateUserRequest) error {

	fmt.Println("------------>", req)

	query := `
		UPDATE users SET 
			id = :id,
			first_name = :first_name,
			last_name = :last_name,
			phone_number = :phone_number
		WHERE id = :id
	
	`
	params := map[string]interface{}{
		"id":           req.Id,
		"first_name":   req.FirstName,
		"last_name":    req.LastName,
		"phone_number": req.PhoneNumber,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	_, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		return nil
	}

	return nil
}

func (u *userRepo) Delete(ctx context.Context, req *user_service.PrimaryKey) (err error) {
	query := `DELETE FROM users WHERE id = $1`

	_, err = u.db.Exec(ctx, query, req.Id)
	if err != nil {
		return err
	}

	return nil
}

// func (u *userRepo) GetUserByPhone(ctx context.Context, req *user_service.PhoneNumber) (resp *user_service.Checker, err error) {

// 	var user user_service.User

// 	fmt.Println("--------------------------RRR", req.PhoneNumber,req)

// 	query := `
// 		SELECT phone_number FROM users WHERE phone_number = $1
// 	`

// 	err = u.db.QueryRow(ctx, query, req.PhoneNumber).Scan(
// 		&user.PhoneNumber,
// 	)
// 	fmt.Println("->>>>-----------> Error has", resp)
// 	if err != nil {
// 		return resp, err
// 	}

// 	return &user_service.Checker{}, err
// }

// func (u *userRepo) GetUserByPhone(ctx context.Context, req *user_service.PrimaryKey) (resp *user_service.Checker, err error) {
// 	var user user_service.User

// 	query := `
// 		SELECT phone_number, name, email FROM users WHERE phone_number = $1
// 	`

// 	err := u.db.QueryRow(ctx, query, req.PhoneNumber).Scan(
// 		&user.PhoneNumber,
// 		&user.Name,
// 		&user.Email,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return &user_service.Checker{
// 				user.Check: false,
// 			}, nil
// 		} else {
// 			return nil, err
// 		}
// 	}

// 	return &user_service.Checker{
// 		Valid: true,
// 	}, nil
// }

func (u *userRepo) Login(ctx context.Context, req *user_service.LoginRequest) (*user_service.LoginResponse, error) {
	phoneNumber := req.PhoneNumber

	var userID int64
	query := `SELECT id FROM users WHERE phone_number = ?`
	err := u.db.QueryRow(ctx, query, phoneNumber).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "User not found")
		}
		return nil, status.Errorf(codes.Internal, "Database error: %v", err)
	}
	return &user_service.LoginResponse{UserId: userID}, nil
}

func (u *userRepo) Register(ctx context.Context, req *user_service.CreateUserRequest) (*user_service.RegisterResponse, error) {
	phoneNumber := req.PhoneNumber
	firstName := req.FirstName
	lastName := req.LastName

	var existingID int64
	query := "SELECT id FROM users WHERE phone_number = ?"
	err := u.db.QueryRow(ctx, query, phoneNumber).Scan(&existingID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, status.Errorf(codes.Internal, "Database error: %v", err)
		}
	} else {
		return nil, err
	}

	query = `INSERT INTO users (phone_number, first_name, last_name) VALUES ($1, $2, $3)`
	result, err := u.db.Exec(ctx, query, phoneNumber, firstName, lastName)
	if err != nil {
		return nil, err
	}

	fmt.Println(result)
	// userID, err := user_service.User.PhoneNumber
	// if err != nil {
	// 	return nil, err
	// }

	// smsReq := &sms_service.GetSmsRequest{SmsId: phoneNumber}
	// if _, err := u.strg.SendSmsCode(ctx, smsReq); err != nil {
	// 	return nil, err
	// }

	return nil, nil
}
