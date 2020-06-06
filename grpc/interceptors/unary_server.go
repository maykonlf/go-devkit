package interceptors

import (
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

func UnaryServerRecovery() grpc.UnaryServerInterceptor {
	return grpc_recovery.UnaryServerInterceptor()
}
