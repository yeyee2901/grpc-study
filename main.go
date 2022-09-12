package main

import (
    // sys
	"fmt"
	"log"
	"os"

	// protobuf stuff
	"google.golang.org/protobuf/proto"

    // gRPC stuffs

	bookpb "yeyee2901/protobuf/gen/go/book/v1"
)

func main() {

	// create protobuf
	someUser := bookpb.Book{
		Title: "Protobuf Tutorial",
		Tahun: 2000,
		Isbn:  "123456789",
	}

	// WRITER ------------------------------------------
	byteBuf, err := proto.Marshal(&someUser)

	if err != nil {
		log.Fatalln("[PROTOBUF] Failed to marshal")
	}

	if err := os.WriteFile("result.bin", byteBuf, 0644); err != nil {
		log.Fatalln("[OS] Failed to binary write to file", err)
	}

	// READER ------------------------------------------
	byteRes, err := os.ReadFile("./result.bin")

	if err != nil {
		log.Fatalln("[OS] Failed to read file", err)
	}

	resultUnmarshal := new(bookpb.Book)

	err = proto.Unmarshal(byteRes, resultUnmarshal)

	if err != nil {
		log.Fatalln("[PROTOBUF]", err)
	}

	fmt.Println(resultUnmarshal)
}
