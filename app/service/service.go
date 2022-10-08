package service

import (
	bookService "yeyee2901/grpc/app/book"
	"yeyee2901/grpc/app/datasource"
	userService "yeyee2901/grpc/app/user"
	bookpb "yeyee2901/grpc/gen/proto/book/v1"
	userpb "yeyee2901/grpc/gen/proto/user/v1"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

type Service struct {
	GRPCServer *grpc.Server
	DB         *sqlx.DB
}

func NewService(grpcServer *grpc.Server, db *sqlx.DB) *Service {
	return &Service{
		GRPCServer: grpcServer,
		DB:         db,
	}
}

func (s *Service) RegisterGRPCServices() {
	// create datasource
	ds := datasource.NewDataSource(s.DB)

	// book handler
	bookServer := bookService.NewBookService(ds)
	bookpb.RegisterBookServiceServer(s.GRPCServer, bookServer)

	// user handler
	userServer := userService.NewUserService(ds)
	userpb.RegisterUserServiceServer(s.GRPCServer, userServer)
}
