package main

import (
	"context"
	"log"
	"net"

	bookpb "yeyee2901/protobuf/gen/go/book/v1"

	"google.golang.org/grpc"
)

type bookService struct{}

func (bs *bookService) GetBook(_ context.Context, req *bookpb.GetBookRequest) (*bookpb.GetBookResponse, error) {
	return &bookpb.GetBookResponse{
		Book: &bookpb.Book{
			Title: req.Title,
			Isbn:  "123456789",
			Tahun: 2022,
		},
	}, nil
}

func main() {

	listener, err := net.Listen("tcp", "localhost:3030")

	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	grpcServer := grpc.NewServer()
	bookpb.RegisterBookServiceServer(grpcServer, &bookService{})

	if err := grpcServer.Serve(listener); err != nil {
		log.Printf("%v", err)
	}
}
