# gRPC Golang
My journey on studying gRPC implementation in golang

# Creating Protobuf
- Install `buf`, this is a wrapper program for `protoc` (protobuf compiler)
```bash
go install github.com/bufbuild/buf/cmd/buf@latest

# plugin for golang template
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```
- Create `buf.yaml`, this contains the options for linting our protobuf files
```yaml
version: v1
breaking:
  use:
    - FILE
lint:
  use:
    - DEFAULT
```
- Create `buf.gen.yaml`, this contains the options controlling the actual protobufs compilation
```yaml
version: v1
plugins:
  - name: go
    out: gen/go                    # output directory prefix
    opt: paths=source_relative     # where the paths are located
  - name: go-grpc
    out: gen/go
    opt: 
      - paths=source_relative
      - require_unimplemented_servers=false # optional
```
- If we want to generate stubs for another language, we can do so by specifying a new entry for the specified language, refer to [documentation](https://docs.buf.build/configuration/v1/buf-gen-yaml)
- Now create the protobuf files, and simply run `buf lint` to lookup for errors
- After creating the protobuf files, run `buf generate` to generate the stubs for your language of choice!

&nbsp;  
  



# Sample Protobuf Service Definition
### /book/v1/book.proto
Contains the model definitions. Make sure the `package` & `go_package` matches what `opt` you specified in `buf.gen.yaml`
```proto
syntax = "proto3";

package book.v1;

option go_package = "yeyee2901/protobuf/gen/go/book/v1;bookpb";

message Book {
  string title = 1;
  string isbn  = 2;
  uint32 tahun = 3;
}
```
&nbsp;


### /book/v1/book_service.proto
Contains the service definition
```proto
syntax = "proto3";

package book.v1;

option go_package = "yeyee2901/protobuf/gen/go/book/v1;bookpb";

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
- Create TCP socket
- create new grpc server object
- register the services
```go
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
