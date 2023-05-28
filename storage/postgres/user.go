package postgres

import (
	"bicycle/bicycle_go_user_service/genproto/user_service"
	"bicycle/bicycle_go_user_service/storage"
	"context"
	"database/sql"

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


// l
// func (u *userRepo) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
//     phoneNumber := req.PhoneNumber

//     // check if user with this phone number exists in database
//     var userID int64
//     query := "SELECT id FROM users WHERE phone_number = $1"
//     err := u.db.QueryRow(query, phoneNumber).Scan(&userID)
//     if err != nil {
//         if err == sql.ErrNoRows {
//             return nil, status.Errorf(codes.NotFound, "User not found")
//         }
//         return nil, status.Errorf(codes.Internal, "Database error: %v", err)
//     }

//     // user found, send success response with user ID
//     return &user_service.LoginResponse{UserId: userID}, nil
// }

func (u *userRepo) Login(ctx context.Context, req *user_service.LoginRequest) (*user.LoginResponse, error) {
    phoneNumber := req.PhoneNumber

    // check if user with this phone number exists in database
    var userID int64
    query := "SELECT id FROM users WHERE phone_number = $1"
    err := u.db.QueryRow(query, phoneNumber).Scan(&userID)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, status.Errorf(codes.NotFound, "User not found")
        }
        return nil, status.Errorf(codes.Internal, "Database error: %v", err)
    }


    return &user.LoginResponse{UserId: userID}, nil
}