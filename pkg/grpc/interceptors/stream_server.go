package interceptors

import (
	"strings"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/maykonlf/go-devkit/pkg/log"
	"google.golang.org/grpc"
)

func StreamServerInterceptors() grpc.ServerOption {
	return grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		StreamServerLogger(),
		StreamServerRecovery(),
	))
}

func StreamServerRecovery() grpc.StreamServerInterceptor {
	return grpc_recovery.StreamServerInterceptor()
}

func StreamServerLogger() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()
		err := handler(srv, stream)
		duration := time.Since(startTime)
		code := grpc_logging.DefaultErrorToCode(err)
		values := strings.Split(info.FullMethod, "/")
		log.Log(getCodeLevel(code), "finished streaming call with code "+code.String(),
			"grpc_service", values[1], "grpc_method", values[2], "grpc_code", code.String(),
			"grpc_start_time", startTime, "grpc_deadline", startTime.Add(duration), "grpc_duration", duration.String(),
		)
		return err
	}
}
