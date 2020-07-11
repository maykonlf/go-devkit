package main

import (
	"errors"
	"github.com/maykonlf/go-devkit/pkg/grpc/server"
)

func main() {
	s := server.NewServer("example-server", 9090)
	s.AddHealthChecks(myFirstHealthCheck, mySecondHealthCheck)
	s.Serve()
}

func myFirstHealthCheck() error {
	return nil
}

func mySecondHealthCheck() error {
	return errors.New("some error")
}
