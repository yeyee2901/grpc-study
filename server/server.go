package server

import (
    "context"

    bookpb "yeyee2901/protobuf/gen/go/book/v1"
)

type bookService struct { } 

func (bs *bookService) GetBook(_ context.Context, req *bookpb.GetBookRequest) (*bookpb.GetBookResponse, error) {
    return &bookpb.GetBookResponse {
        Book: &bookpb.Book {
            Title: req.Title,
            Isbn: "123456789",
            Tahun: 2022,
        },
    }, nil
}
