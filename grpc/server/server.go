package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/maykonlf/go-devkit/grpc/interceptors"
	"github.com/maykonlf/go-devkit/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	ctx        context.Context
	name       string
	mux        *runtime.ServeMux
	gRPCServer *grpc.Server
}

func NewServer(ctx context.Context, name string) *Server {
	gRPCServer := grpc.NewServer(interceptors.UnaryServerInterceptors())
	return &Server{
		ctx:        ctx,
		name:       name,
		mux:        runtime.NewServeMux(),
		gRPCServer: gRPCServer,
	}
}

func (s *Server) EnableGRPCReflection() *Server {
	reflection.Register(s.gRPCServer)
	return s
}

func (s *Server) RegisterGRPCServerHandlersFunc(registerFunc func(server *grpc.Server)) {
	registerFunc(s.gRPCServer)
}

func (s *Server) RegisterHTTPServerHandlersFunc(registerFunc func(ctx context.Context, mux *runtime.ServeMux)) {
	registerFunc(s.ctx, s.mux)
}

func (s *Server) ListenAndServe(addr string) error {
	log.Info(fmt.Sprintf("server listening on %s...", addr))
	return http.ListenAndServe(addr, s.getMuxHandler())
}

func (s *Server) getMuxHandler() http.Handler {
	return h2c.NewHandler(http.HandlerFunc(s.httpGRPCMuxHandler), &http2.Server{})
}

func (s *Server) httpGRPCMuxHandler(w http.ResponseWriter, r *http.Request) {
	if isGRPCRequest(r) {
		s.gRPCServer.ServeHTTP(w, r)
	} else {
		s.mux.ServeHTTP(w, r)
	}
}

func isGRPCRequest(r *http.Request) bool {
	return r.ProtoMajor == 2 && strings.EqualFold(r.Header.Get("Content-Type"), "application/grpc")
}
