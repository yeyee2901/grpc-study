package main

import (
	"context"
	"fmt"
	"log"
	bookpb "yeyee2901/grpc/gen/book/v1"

	"google.golang.org/grpc"
)

func main() {
	// for connection without creds, we have to use this
	var opts = []grpc.DialOption{
		grpc.WithInsecure(),
	}

	// bind the grpc connection to localhost:3030 TCP socket
	conn, err := grpc.Dial("localhost:3030", opts...)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// create book client
	bookClient := bookpb.NewBookServiceClient(conn)

	// try requesting
	req := &bookpb.GetBookRequest{
		Title: "Book #1",
	}

	bookResp, err := bookClient.GetBook(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s - %s - %d\n",
		bookResp.Book.GetTitle(),
		bookResp.Book.GetIsbn(),
		bookResp.Book.GetTahun(),
	)
}
