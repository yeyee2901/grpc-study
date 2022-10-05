package main

import (
	"context"
	"fmt"
	"log"
	"net"

	bookpb "yeyee2901/protobuf/gen/go/book/v1"

	"google.golang.org/grpc"
)

type bookService struct{}

func (bs *bookService) GetBook(_ context.Context, req *bookpb.GetBookRequest) (*bookpb.GetBookResponse, error) {
    fmt.Println("Received: ", req.String())
	return &bookpb.GetBookResponse{
		Book: &bookpb.Book{
			Title: req.Title,
			Isbn:  "123456789",
			Tahun: 2022,
		},
	}, nil
}

func main() {

	// create TCP socket
	listener, err := net.Listen("tcp", "localhost:3030")

	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	// create grpc server
	grpcServer := grpc.NewServer()

	// register services
	bookServer := new(bookService)
	bookpb.RegisterBookServiceServer(grpcServer, bookServer)

	// bind the grpc server in the tcp socket
	if err := grpcServer.Serve(listener); err != nil {
		log.Printf("%v", err)
	}
}
