package book

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"yeyee2901/grpc/app/datasource"
	bookpb "yeyee2901/grpc/gen/proto/book/v1"
)

type BookService struct {
	DataSource *datasource.DataSource
}

func NewBookService(ds *datasource.DataSource) *BookService {
	return &BookService{
		DataSource: ds,
	}
}

func (bs *BookService) GetBook(_ context.Context, req *bookpb.GetBookRequest) (*bookpb.GetBookResponse, error) {
	// return the result
	newBook := &bookpb.Book{
		Title: req.Title,
		Isbn:  "123456789",
		Tahun: 2022,
	}

	resp := &bookpb.GetBookResponse{
		Book: newBook,
	}

	return resp, nil
}

func (T *BookService) GetBookStream(req *bookpb.GetBookRequest, stream bookpb.BookService_GetBookStreamServer) (err error) {
	books := generateDummyBooks(5)

	for _, book := range books {
		// simulate waiting
		time.Sleep(time.Duration(1) * time.Second)

		err := stream.Send(&bookpb.GetBookResponse{
			Book: book,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (bs *BookService) SaveBook(_ context.Context, req *bookpb.SaveBookRequest) (resp *bookpb.SaveBookResponse, err error) {
	book := req.Book

	err = bs.DataSource.SaveBook(book)

	if err != nil {
        resp = &bookpb.SaveBookResponse{
        	Message:   "[ERROR] "+err.Error(),
        	NewBookId: 0,
        }
		return nil, err
	}

	resp = &bookpb.SaveBookResponse{
		Message:   "Sukses. Abaikan NewBookId, karena masih generate random, masih ngoprek RETURNING clause di driver pgx - sqlx golang",
		NewBookId: int64(rand.Int()),
	}

	return
}

func generateDummyBooks(count int) []*bookpb.Book {
	var books []*bookpb.Book

	for i := 0; i < count; i++ {
		title := fmt.Sprintf("Book #%d", i)
		isbn := fmt.Sprintf("ISBN - %d", i)
		books = append(books, &bookpb.Book{
			Title: title,
			Isbn:  isbn,
			Tahun: uint32(2000 + i),
		})
	}

	return books
}
