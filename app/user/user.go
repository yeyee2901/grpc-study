package user

import (
	"context"
	"database/sql"
	"yeyee2901/grpc/app/datasource"
	userpb "yeyee2901/grpc/gen/proto/user/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	DataSource *datasource.DataSource
}

func (us *UserService) GetUserById(_ context.Context, req *userpb.GetUserByIdRequest) (*userpb.GetUserByIdResponse, error) {
	var (
		user userpb.User
		resp *userpb.GetUserByIdResponse
	)

	// get dari postgre DB
	err := us.DataSource.GetUserById(&user, req.GetId())

	// user not found
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.NotFound, "[NOT FOUND]: User not found")
	}

	// error lain
	if err != nil {
		return nil, status.Error(codes.Internal, "Database Error. "+err.Error())
	}

	// construct response message
	resp = &userpb.GetUserByIdResponse{
		User: &user,
	}

	return resp, nil
}

func (us *UserService) CreateUser(_ context.Context, _ *userpb.CreateUserRequest) (resp *userpb.CreateUserResponse, err error) {
	return
}

func NewUserService(ds *datasource.DataSource) *UserService {
	return &UserService{
		DataSource: ds,
	}
}
