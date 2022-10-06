package service

import (
	bookService "yeyee2901/grpc/app/book"
	bookpb "yeyee2901/grpc/gen/book/v1"

	"google.golang.org/grpc"
)

type Service struct {
	GRPCServer *grpc.Server
}

func NewService(grpcServer *grpc.Server) *Service {
	return &Service{
		GRPCServer: grpcServer,
	}
}

func (s *Service) RegisterGRPCServices() {
	// book handler
	bookServer := bookService.NewBookService()
	bookpb.RegisterBookServiceServer(s.GRPCServer, bookServer)
}
