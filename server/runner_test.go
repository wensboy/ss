package server

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type MockServer struct{}

func (ms *MockServer) Start() {
	fmt.Println("start server...")
}

func (ms *MockServer) GoStart(ctx context.Context) {
	fmt.Println("go start server...")
}

func (ms *MockServer) Stop() {
	fmt.Println("stop server...")
}

var ms = &MockServer{}

func Test_Runner(t *testing.T) {
	runner := NewRunner().SetExitDelay(500 * time.Millisecond).SetServer(ms)
	runner.Run() // block here! ctrl + c to exit...
}

func Test_Runner_RestServer(t *testing.T) {
	s := NewRestServer()
	s.MountModules(
		s.MountDatabases(nil, func(s *RestServer) {
			fmt.Println("mount databases")
		}, nil),
		s.MountConfig(nil, func(s *RestServer) {
			fmt.Println("mount config")
			s.server.Addr = ":18080"
		}, nil),
	)
	runner := NewRunner().SetExitDelay(500 * time.Millisecond).SetServer(s)
	runner.Run() // block here! ctrl + c to exit...
}
