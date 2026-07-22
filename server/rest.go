package server

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type RestServer struct {
	muxer  *echo.Echo
	server *http.Server
}

func (s *RestServer) Start() {}

func (s *RestServer) Stop() {}

func (s *RestServer) GoStart() {}

func (s *RestServer) SetMuxer(muxer *echo.Echo) *RestServer {
	s.muxer = muxer
	return s
}

func (s *RestServer) SetServer(server *http.Server) *RestServer {
	s.server = server
	return s
}

func NewRestServer() *RestServer {
	return &RestServer{
		muxer:  echo.New(),
		server: &http.Server{},
	}
}
