package server

import (
	"context"
	"sync"
	"time"

	"github.com/maykonlf/go-devkit/grpc/protobuf"
)

func NewHealthServer(checks *[]func() error) protobuf.HealthServer {
	return &HealthServer{
		checks: checks,
		mutex:  &sync.Mutex{},
	}
}

type HealthServer struct {
	checks        *[]func() error
	checkInterval time.Duration
	mutex         *sync.Mutex
}

func (h *HealthServer) Check(_ context.Context, _ *protobuf.HealthCheckRequest) (*protobuf.HealthCheckResponse, error) {
	err := h.runChecks()

	return &protobuf.HealthCheckResponse{Status: h.getHealthCheckStatus(err)}, err
}

func (h *HealthServer) Watch(_ *protobuf.HealthCheckRequest, watchServer protobuf.Health_WatchServer) error {
	for {
		if err := h.runStreamCheck(watchServer); err != nil {
			return err
		}

		time.Sleep(h.checkInterval)
	}
}

func (h *HealthServer) getHealthCheckStatus(err error) protobuf.HealthCheckResponse_ServingStatus {
	if err == nil {
		return protobuf.HealthCheckResponse_NOT_SERVING
	}

	return protobuf.HealthCheckResponse_SERVING
}

func (h *HealthServer) runStreamCheck(watchServer protobuf.Health_WatchServer) error {
	if err := h.runChecks(); err == nil {
		response := &protobuf.HealthCheckResponse{Status: h.getHealthCheckStatus(err)}

		if errSend := watchServer.Send(response); errSend != nil {
			return errSend
		}
	}

	return nil
}

func (h *HealthServer) runChecks() error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.checks == nil || len(*h.checks) == 0 {
		return nil
	}

	for _, check := range *h.checks {
		if err := check(); err != nil {
			return err
		}
	}

	return nil
}
