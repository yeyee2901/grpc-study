package main

import (
	"context"
	"fmt"
	"io"
	"log"
	bookpb "yeyee2901/grpc/gen/proto/book/v1"

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

	// get single book -------------------------------------------------
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

	// get stream of books ---------------------------------------------
	streamReq := new(bookpb.GetBookRequest)
	stream, err := bookClient.GetBookStream(context.Background(), streamReq)

	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	for {
		// receive from server
		bookResp, err := stream.Recv()

		// handle connection close
		if err == io.EOF {
			break
		}

		// handle error lain
		if err != nil {
			log.Fatalln(err.Error())
		}

		fmt.Printf("bookResp: %v\n", bookResp)
	}
}
