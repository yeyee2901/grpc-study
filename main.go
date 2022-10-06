package main

import (
	"log"
	"net"
	"net/url"

	"yeyee2901/grpc/app/service"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

type App struct {
	Listener   net.Listener
	GRPCServer *grpc.Server
	DB         *sqlx.DB
}

func (app *App) InitGRPCServices() {
	service := service.NewService(app.GRPCServer, app.DB)
	service.RegisterGRPCServices()
}

func (app *App) InitDatabase() {
	dsConn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword("your_username", "your_password"),
		Host:   "localhost:5432",
		Path:   "local_development",
	}

	// disable ssl since its only in localhost
	optsQuery := dsConn.Query()
	optsQuery.Add("sslmode", "disable")
	dsConn.RawQuery = optsQuery.Encode()

	// open the connection
	db, err := sqlx.Open("pgx", dsConn.String())

	if err != nil || db == nil {
		log.Fatalln("Failed to connect to database", err)
	}

	app.DB = db
}

func main() {
	app := App{}

	// INIT: Create TCP socket
	listener, err := net.Listen("tcp", "localhost:3030")

	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	app.Listener = listener

	// INIT: databse PostgreSQL
	app.InitDatabase()

	// TODO: INIT interceptors (middlewares)

	// INIT: gRPC services
	grpcServer := grpc.NewServer()
	app.GRPCServer = grpcServer
	app.InitGRPCServices()

	// Start the server :)
	if err := app.GRPCServer.Serve(listener); err != nil {
		log.Printf("%v", err)
	}
}
