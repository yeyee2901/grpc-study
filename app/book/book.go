package book

import (
	"context"
	"fmt"

	bookpb "yeyee2901/grpc/gen/book/v1"

	"google.golang.org/grpc"
)

type BookService struct {
	GRPCServer *grpc.Server
}

func NewBookService() *BookService {
	return &BookService{}
}

func (bs *BookService) GetBook(_ context.Context, req *bookpb.GetBookRequest) (*bookpb.GetBookResponse, error) {
	// print the request
	fmt.Println("Received: ", req.String())

	// return the result
	resp := &bookpb.GetBookResponse{
		Book: &bookpb.Book{
			Title: req.Title,
			Isbn:  "123456789",
			Tahun: 2022,
		},
	}
	return resp, nil
}
