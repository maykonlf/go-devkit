package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/maykonlf/go-devkit/pkg/grpc/protobuf"
	"github.com/maykonlf/go-devkit/pkg/grpc/server"
	"google.golang.org/grpc"
	"net/http"
)

func runGRPCServer() {
	s := server.NewServer("example-server", 9090)
	s.AddHealthChecks(func() error {
		return nil
	})

	panic(s.Serve())
}

func runGRPCGateway() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := protobuf.RegisterHealthHandlerFromEndpoint(ctx, mux, "localhost:9090", opts)
	if err != nil {
		panic(err)
	}

	panic(http.ListenAndServe(":8080", mux))
}

func main() {
	go runGRPCServer()
	runGRPCGateway()
}
