package main

import (
	"log"
	"net"

	"yeyee2901/grpc/app/service"

	"google.golang.org/grpc"
)

type App struct {
	Listener   net.Listener
	GRPCServer *grpc.Server
}

func (app *App) InitGRPCServices() {
	service := service.NewService(app.GRPCServer)
	service.RegisterGRPCServices()
}

func main() {
  app := App{}

  // INIT: Create TCP socket
	listener, err := net.Listen("tcp", "localhost:3030")

	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

  app.Listener = listener

  // INIT: gRPC services
	grpcServer := grpc.NewServer()
  app.GRPCServer = grpcServer
  app.InitGRPCServices()

	// bind the grpc server in the tcp socket
	if err := app.GRPCServer.Serve(listener); err != nil {
		log.Printf("%v", err)
	}
}
