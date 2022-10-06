package interceptors

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// LoggerInterceptor godoc
// @Summary logs request & response to log file. Pass this to
func LoggerUnaryRPC(
	ctx context.Context, req interface{}, serverInfo *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {

	// log time
	start := time.Now()

	// Log request
	log.Info().
		Str("method", serverInfo.FullMethod).
		Str("req", fmt.Sprintf("%+v", req)).
		Msg("REQUEST")

	// call the RPC handler (equivalent to ctx.Next() in gin)
	resp, err = handler(ctx, req)

	// after RPC handler invoked
	latency := time.Since(start).Seconds()
	log.Info().
		Str("method", serverInfo.FullMethod).
		Str("req", fmt.Sprintf("%+v", req)).
		Str("resp", fmt.Sprintf("%+v", resp)).
		Float32("latency", float32(latency)).
		Msg("RESPONSE")

	return
}
