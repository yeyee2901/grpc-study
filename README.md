# gRPC Golang
My journey on studying gRPC implementation in golang  
&nbsp;  
**Target**:
- ✅ Integrate with `postgre` database (using `sqlx` package)
- Create CRUD gRPC API
- Try out protobuf for react-typescript (for front end) >> might be in separate repository

# UPDATES
- `2022-10-07`: Added unary RPC logger interceptor using [zerolog](https://github.com/rs/zerolog). You can monitor the log file using `tail -f log/service.log`
- `2022-10-06`: Added [protoc-go-inject-tag](https://github.com/favadi/protoc-go-inject-tag) for custom struct tags injection, mainly used for easier integration with database layer using [jmoiron/sqlx](https://github.com/jmoiron/sqlx)

# Project Structure (Best Practice from Google)
```
root
├── app
│   ├── service                                 # Register gRPC services here
│   │   └── service.go
│   ├── interceptors                            # Interceptors are defined here, but are registered in main when initiating the server
│   │   └── interceptors.go                     # this is equivalent to middlewares in HTTP
│   ├── datasource                              # datasource, for now it's only postgre database
│   │   └── datasource.go
│   └── book                                    # book services stubs & logics are implemented here
│       └── book.go
├── client                                      # for demonstration purposes
│   └── client.go
├── gen                                         # generated stubs
│   └── proto
│       └── book                                # book stubs
│           └── v1
│               ├── book.pb.go
│               ├── book_service_grpc.pb.go
│               └── book_service.pb.go
├── proto                                       # proto files
│   └── book                                    # proto files related to books
│       └── v1
│           ├── book.proto
│           └── book_service.proto
├── go.mod
├── go.sum
├── schema.sql
├── main.go                                     # entry point
├── Makefile                                    # for simplifying build command
└── README.md                                   # readme senpai~
```

# How to Try This Project
1. Clone the repository
2. Download the dependencies
```bash
go mod download && go mod verify
```
3. Run it the server
```bash
go run .
```
4. Open your Postman (yes, as of Postman 10, it now supports gRPC calls)
5. Create a new workspace
6. Create a new request
![image](https://user-images.githubusercontent.com/55247343/194440939-0fc35363-1e80-41bd-a478-9356edbf3a8f.png)
![image](https://user-images.githubusercontent.com/55247343/194440981-1a9c4fc6-95c7-4c5f-b57a-9002f87d8866.png)
7. Create a new API definition using import method
![image](https://user-images.githubusercontent.com/55247343/194441086-2b95a9b5-089b-422c-bc86-e89904af804f.png)
8. Navigate it to service proto files (the one that ends with `_service.proto`). We do this since the service proto files already imported the model definition anyways, and also the RPC definitions are stored in the `_service.proto`, so this is the one that Postman needs to import.
9. Add import path, this is crucial so the `_service.proto` can resolve its import statements. Navigate it to the **root project directory**. After you're done, click **Next**.
![image](https://user-images.githubusercontent.com/55247343/194441328-a8e91701-9e5b-49c5-8946-c8f017cdc4f1.png)
10. Type in your API name (any name you want), after you're done, click **Import as API**.
![image](https://user-images.githubusercontent.com/55247343/194441512-e942bdf9-e83c-4a54-a264-341ab638792c.png)
11. Set the server URL to `localhost:3030` (project default)
12. Done! now you can 'invoke' the RPC methods!
![image](https://user-images.githubusercontent.com/55247343/194441590-305cbc77-a386-446c-9de8-eb4576937d6a.png)
13. Note that you can also generate sample protobuf message using Postman. Postman is smart enough to recognize the request message definition.
![image](https://user-images.githubusercontent.com/55247343/194441731-d52a78dd-7586-4dbe-9673-e93d35511182.png)
14. Also, if you change any of the proto files, you will need to re import it in the Postman 😜

# Installing Tools for Working with Protobuf
We need 3 things:
- `protoc`, the protobuf compiler. This is responsible for generating the code & stubs for your language of choice. Check your system for how to install these. I used Ubuntu & Arch, so I can provide the package name, otherwise you might need to install from [source]( https://developers.google.com/protocol-buffers/)
- `protoc-gen-go`, a plugin for `protoc` which provides the capability of generating `golang` interfaces
- `protoc-gen-go-grpc`, same as `protoc-gen-go`, but also provides the interface for integrating with gRPC services
- `protoc-go-inject-tag`, this is useful for injecting custom tags. I use this because `sqlx` has a custom tag that allows object binding. Otherwise, you might need to "inject" your `.pb.go` file with your custom tags everytime you change your proto files, which is a hassle. This program simplifies that "tags injection" process.
```bash
# arch based system available in AUR
yay -S protobuf

# ubuntu
sudo apt install protobuf-compiler

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/favadi/protoc-go-inject-tag@latest
```
- Create `Makefile`, we will define our protobuf compile command & options here
```make
BASE_OUTPUT_DIR = ./gen


all: proto-go

proto-go:
	protoc  --go_out=${BASE_OUTPUT_DIR} --go_opt=paths=source_relative \
			--go-grpc_out=${BASE_OUTPUT_DIR} --go-grpc_opt=paths=source_relative \ 
			--go-grpc_opt=require_unimplemented_servers=false \
			book/v1/*.proto
	protoc-go-inject-tag -input="./gen/proto/book/v1/*.pb.go"
```
*Note: the `require_unimplemented_servers=false` will disable forward compatibility. For most cases, this shouldn't matter much.
&nbsp;  
  



# Sample Protobuf Service Definition
### /book/v1/book.proto
Contains the model definitions.
```proto
syntax = "proto3";

package book.v1;

option go_package = "yeyee2901/grpc/gen/book/v1;bookpb";

message Book {
  // @gotags: db:"title"
  string title = 1;

  // @gotags: db:"isbn"
  string isbn = 2;

  // @gotags: db:"tahun"
  uint32 tahun = 3;
}

```
&nbsp;


### /book/v1/book_service.proto
Contains the service definition
```proto
syntax = "proto3";

package book.v1;

option go_package = "yeyee2901/protobuf/gen/book/v1;bookpb";

import "book/v1/book.proto";

service BookService {
  rpc GetBook(GetBookRequest) returns (GetBookResponse) {};
}

message GetBookRequest {
  string title = 1;
}

message GetBookResponse {
  Book book = 1;
}
```
&nbsp;




# Creating gRPC Server
- Implement the `GetBook` method from `BookServiceServer`
- Create TCP socket
- create new grpc server object
- register the services
```go
type bookService struct{}

// implement the GetBook method, this is required so we can register the service to the gRPC server
func (bs *bookService) GetBook(_ context.Context, req *bookpb.GetBookRequest) (*bookpb.GetBookResponse, error) {
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


// in main .......
// create TCP socket
listener, err := net.Listen("tcp", "localhost:3030")

if err != nil {
	log.Fatalf("Failed to listen %v", err)
}

grpcServer := grpc.NewServer()
// file book/v1/book.proto
// register the services
bookServer := new(bookService)
bookpb.RegisterBookServiceServer(grpcServer, bookServer)

// make grpc server listen on the TCP socket
if err := grpcServer.Serve(listener); err != nil {
	log.Printf("%v", err)
}
```
&nbsp;




# Creating gRPC Clients
- create option object, by default since it's in localhost, it's safe to say we don't need any credentials, so we can use `grpc.WithInsecure()` as the option
- listen on server address, pass in the options (using spread operator)
- create new gRPC client, this client is associated with a single gRPC service type. If we have many grpc service (routing), we have to create clients for each gRPC service that we want to call.
- Create a new request object, it is as defined in the `.proto` file
```go
// for connection without creds, we have to use this
var opts  = []grpc.DialOption{
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
```
&nbsp;




# gRPC Service Types
### 1. Simple RPC
- simple request-response flow, like RESTful API
- client send 1 request
- server send 1 response
- example is on top

&nbsp;  
  
  
  
  
### 2. Server side streaming
- client send 1 request
- server sends many response data

#### Server side code
```go
for {
    if err := stream.Send(data); err != nil {
        return err
    }
    
    if stopCondition { 
        break 
    }
}

return nil  // closes the connection
```
#### Client side code
```go
for {
    respData, err := stream.Recv()
    
    // error pertama yang di cek harus connection close
    if err == io.EOF {
        break
    }
    
    // error lain
    if err != nil {
        log.Fatalln(err)
    }
    
    // handle data
    fmt.Println(respData)
}
```

&nbsp;  
  
  
  
### 3. Client side streaming
- client sends many request data
- server sends 1 response data

#### Server side code
```go
for {
    var dataArray []RequestDataStruct
    data, err := stream.Recv()
    
    // handle error pertama harus cek connection close
    if err == io.EOF {
        return stream.SendAndClose(&ResponseDataStruct{
            Datas: dataArray,
        })
    }
    
    // error lain
    if err != nil {
        return err
    }
    
    // do something with the data...
    dataArray = append(dataArray, data)
}
```
#### Client side Code
```go
for _, data := range dataArray {
    if err := stream.Send(data); err != nil {
        log.Fatalln("Failed to send data", err)
    }
}

// tunggu respon & tutup koneksi
respData, err := stream.CloseAndRecv()
if err != nil {
    // ......
}

```
&nbsp;
  
  
### 4. Bidirectional streaming
- both client & server sends many data until both closes the connection

```go
stream, err := client.RouteChat(context.Background())
done := make(chan struct{})

// for handling server streaming
go func() {
  for {
    in, err := stream.Recv()
    if err == io.EOF {
      // read done.
      close(done)
      return
    }
    if err != nil {
      log.Fatalf("Failed to receive a note : %v", err)
    }
    log.Printf("Got message %s at point(%d, %d)", in.Message, in.Location.Latitude, in.Location.Longitude)
  }
}()

// for sending to server
for _, note := range notes {
  if err := stream.Send(note); err != nil {
    log.Fatalf("Failed to send a note: %v", err)
  }
}

// done? close the connection, on the other side, this will produce err == io.EOF
stream.CloseSend()
<-done
```
