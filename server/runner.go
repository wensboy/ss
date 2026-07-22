package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server interface {
	Start()
	Stop()
	GoStart(context.Context) // goroutine start
}

type Runner struct {
	server    Server
	exitDelay time.Duration
}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) SetServer(s Server) *Runner {
	r.server = s
	return r
}

func (r *Runner) SetExitDelay(delay time.Duration) *Runner {
	r.exitDelay = delay
	return r
}

func (r *Runner) Run() {
	go func() {
		r.server.Start()
	}()
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-exit
	r.Exit()
}

func (r *Runner) Exit() {
	time.Sleep(r.exitDelay)
	r.server.Stop()
}
