package interceptors

import (
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

func StreamServerRecovery() grpc.StreamServerInterceptor {
	return grpc_recovery.StreamServerInterceptor()
}
