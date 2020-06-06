package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/maykonlf/go-devkit/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net/http"
	"strings"
)

type Server struct {
	ctx        context.Context
	name       string
	logger     log.LoggerI
	mux        *runtime.ServeMux
	gRPCServer *grpc.Server
}

func NewServer(ctx context.Context, name string, logger log.LoggerI) *Server {
	return &Server{
		ctx:        ctx,
		name:       name,
		logger:     logger,
		mux:        runtime.NewServeMux(),
		gRPCServer: grpc.NewServer(),
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
	s.logger.Info(fmt.Sprintf("server listening on %s...", addr))
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
