package server

import (
	"fmt"
	"github.com/maykonlf/go-devkit/grpc/interceptors"
	"github.com/maykonlf/go-devkit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Server struct {
	name       string
	port       int
	gRPCServer *grpc.Server
}

func NewServer(name string, port int) *Server {
	gRPCServer := grpc.NewServer(interceptors.UnaryServerInterceptors(), interceptors.StreamServerInterceptors())
	return &Server{
		name:       name,
		port:       port,
		gRPCServer: gRPCServer,
	}
}

func (s *Server) EnableGRPCReflection() *Server {
	reflection.Register(s.gRPCServer)
	return s
}

func (s *Server) Serve() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Panic("failed to initialize listener", "port", s.port, "error", err)
	}

	log.Panic("failed to serve", "error", s.gRPCServer.Serve(lis))
}

func (s *Server) GetGRPCServer() *grpc.Server {
	return s.gRPCServer
}
