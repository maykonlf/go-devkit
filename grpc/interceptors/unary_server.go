package interceptors

import (
	"context"
	"strings"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/maykonlf/go-devkit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var codeLogLevelMap = map[codes.Code]log.Level{
	codes.OK:                 log.InfoLevel,
	codes.Canceled:           log.InfoLevel,
	codes.Unknown:            log.ErrorLevel,
	codes.InvalidArgument:    log.InfoLevel,
	codes.DeadlineExceeded:   log.WarnLevel,
	codes.NotFound:           log.InfoLevel,
	codes.AlreadyExists:      log.InfoLevel,
	codes.PermissionDenied:   log.WarnLevel,
	codes.Unauthenticated:    log.InfoLevel,
	codes.ResourceExhausted:  log.WarnLevel,
	codes.FailedPrecondition: log.WarnLevel,
	codes.Aborted:            log.WarnLevel,
	codes.OutOfRange:         log.WarnLevel,
	codes.Unimplemented:      log.ErrorLevel,
	codes.Internal:           log.ErrorLevel,
	codes.Unavailable:        log.WarnLevel,
	codes.DataLoss:           log.ErrorLevel,
}

func UnaryServerInterceptors() grpc.ServerOption {
	return grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		UnaryServerLogger(),
		UnaryServerRecovery(),
	))
}

func UnaryServerRecovery() grpc.UnaryServerInterceptor {
	return grpc_recovery.UnaryServerInterceptor()
}

func UnaryServerLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		response, err := handler(ctx, req)
		logUnaryRequest(req, info, startTime, response, err)

		return response, err
	}
}

func logUnaryRequest(req interface{},
	info *grpc.UnaryServerInfo,
	startTime time.Time,
	response interface{}, err error) {
	duration := time.Since(startTime)
	code := grpc_logging.DefaultErrorToCode(err)
	values := strings.Split(info.FullMethod, "/")
	log.Log(getCodeLevel(code), "finished unary call with code with "+code.String(),
		"grpc_service", values[1], "grpc_method", values[2], "grpc_code", code.String(),
		"grpc_request", req, "grpc_response", response, "grpc_start_time", startTime,
		"grpc_deadline", startTime.Add(duration), "grpc_duration", duration.String(),
	)
}

func getCodeLevel(code codes.Code) log.Level {
	if level, ok := codeLogLevelMap[code]; ok {
		return level
	}
	return log.ErrorLevel
}
