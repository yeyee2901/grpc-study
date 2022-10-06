package main

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	"yeyee2901/grpc/app/interceptors"
	"yeyee2901/grpc/app/service"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	TIMESTAMP_FORMAT = time.RFC3339
)

type App struct {
	Listener   net.Listener
	GRPCServer *grpc.Server
	DB         *sqlx.DB
}

func (app *App) InitGRPCServices() {
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
		log.Error().Msg("Failed to connect to database")
	}

	app.DB = db
}

func (app *App) InitGRPCServer() {
	serverOpts := []grpc.ServerOption{
		// middlewares
		grpc.UnaryInterceptor(interceptors.LoggerUnaryRPC),
	}

	// init the server with options
	grpcServer := grpc.NewServer(serverOpts...)
	app.GRPCServer = grpcServer

	// register the services
	service := service.NewService(app.GRPCServer, app.DB)
	service.RegisterGRPCServices()
}

func (app *App) InitLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = TIMESTAMP_FORMAT
	log.Logger = zerolog.New(&lumberjack.Logger{
		Filename:   "log/service.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	})
	log.Logger = log.With().Str("service", "gRPC Golang").Logger()
	log.Logger = log.With().Caller().Logger()
	log.Logger = log.With().Timestamp().Logger()
}

func main() {
	app := App{}

	// INIT: Create TCP socket
	listener, err := net.Listen("tcp", "localhost:3030")
	if err != nil {
		fmt.Println("[ERROR] Failed to create listener")
		os.Exit(1)
	}

	app.Listener = listener

	// INIT: databse PostgreSQL
	app.InitDatabase()

	// TODO: INIT custom logger
	app.InitLogger()

	// INIT: gRPC Server
	app.InitGRPCServer()

	// Start the server :)
	fmt.Println("Server starting ....")
	if err := app.GRPCServer.Serve(listener); err != nil {
		fmt.Println("Server exited ...")
		os.Exit(1)
	}
}
