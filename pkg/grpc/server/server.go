package server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/maykonlf/go-devkit/pkg/grpc/interceptors"
	"github.com/maykonlf/go-devkit/pkg/grpc/protobuf"
	"github.com/maykonlf/go-devkit/pkg/log"
)

type Server struct {
	name         string
	port         int
	gRPCServer   *grpc.Server
	healthChecks []func() error
}

func NewServer(name string, port int) *Server {
	gRPCServer := grpc.NewServer(interceptors.UnaryServerInterceptors(), interceptors.StreamServerInterceptors())
	server := &Server{
		name:         name,
		port:         port,
		gRPCServer:   gRPCServer,
		healthChecks: []func() error{},
	}

	protobuf.RegisterHealthServer(server.gRPCServer, NewHealthServer(&server.healthChecks))

	return server
}

func (s *Server) EnableGRPCReflection() *Server {
	reflection.Register(s.gRPCServer)
	return s
}

func (s *Server) AddHealthChecks(checks ...func() error) *Server {
	s.healthChecks = append(s.healthChecks, checks...)
	return s
}

func (s *Server) GetGRPCServer() *grpc.Server {
	return s.gRPCServer
}

func (s *Server) Serve() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Panic("failed to initialize listener", "port", s.port, "error", err)
	}

	log.Panic("failed to serve", "error", s.gRPCServer.Serve(lis))
}
