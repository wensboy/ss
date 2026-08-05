package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/wensboy/ss/config"
	"github.com/wensboy/ss/err"
)

const (
	GLOBAL_MIDDLEWARE = "_global_middleware"
)

type RestServer struct {
	muxer    *echo.Echo
	server   *http.Server
	routers  map[string]*echo.Group // 用于中间件追踪挂载
	Scontext *ServerContext
}

func NewRestServer() *RestServer {
	return &RestServer{
		muxer:    echo.New(),
		server:   &http.Server{},
		routers:  make(map[string]*echo.Group),
		Scontext: NewServerContext(),
	}
}

type RestServerOption func(*RestServer)
type RestServerHook = RestServerOption

func (s *RestServer) Start() {
	// 确保执行过程中一定有 server
	s.server.Handler = s.muxer
	if perr := s.server.ListenAndServe(); perr != nil && perr != http.ErrServerClosed {
		panic(err.GetGErrHub().GetOrSet(err.NewNormalErrCode(ErrCodeStartFail), err.NewErr("", err.NewNormalErrCode(ErrCodeStartFail).Code(), "server start failed").Wrap(perr)))
	}
}

func (s *RestServer) Stop() {
	if perr := s.server.Shutdown(context.Background()); perr != nil {
		panic(err.GetGErrHub().GetOrSet(err.NewNormalErrCode(ErrCodeStopFail), err.NewErr("", err.NewNormalErrCode(ErrCodeStopFail).Code(), "server stop failed").Wrap(perr)))
	}
}

func (s *RestServer) GoStart(ctx context.Context) {
	go func() {
		s.Start()
	}()
	go func() {
		<-ctx.Done()
		s.Stop()
	}()
}

func (s *RestServer) SetMuxer(muxer *echo.Echo) *RestServer {
	s.muxer = muxer
	return s
}

func (s *RestServer) SetServer(server *http.Server) *RestServer {
	s.server = server
	return s
}

func (s *RestServer) SetServerContext(scontext *ServerContext) *RestServer {
	s.Scontext = scontext
	return s
}

func (s *RestServer) MountModules(opts ...RestServerHook) {
	for _, opt := range opts {
		opt(s)
	}
}

func (s *RestServer) MountRouters(prefix string, routers ...func(*echo.Group, *ServerContext)) RestServerOption {
	if s.routers[prefix] == nil {
		s.routers[prefix] = s.muxer.Group(prefix)
	}
	return func(s *RestServer) {
		for _, router := range routers {
			router(s.routers[prefix], s.Scontext)
		}
	}
}

// 总是在 router 之前挂载中间件
func (s *RestServer) MountMiddlewares(prefix string, middlewares ...echo.MiddlewareFunc) RestServerOption {
	if prefix != GLOBAL_MIDDLEWARE && s.routers[prefix] == nil {
		s.routers[prefix] = s.muxer.Group(prefix)
	}
	return func(s *RestServer) {
		if prefix == GLOBAL_MIDDLEWARE {
			s.muxer.Use(middlewares...)
			return
		}
		s.routers[prefix].Use(middlewares...)
	}
}

func (s *RestServer) MountDatabases(preHook, runHook, postHook RestServerHook) RestServerOption {
	return func(s *RestServer) {
		if preHook != nil {
			preHook(s)
		}
		if runHook != nil {
			runHook(s)
		} else {
			s.Scontext.MountDBContext()
		}
		if postHook != nil {
			postHook(s)
		}
	}
}

func (s *RestServer) MountConfig(preHook, runHook, postHook RestServerHook) RestServerOption {
	return func(s *RestServer) {
		if preHook != nil {
			preHook(s)
		}
		if runHook != nil {
			runHook(s)
		} else {
			s.mountConfig()
		}
		if postHook != nil {
			postHook(s)
		}
	}
}

func (s *RestServer) mountConfig() {
	s.server.Addr = config.MustLookup[string](
		config.GFlagSource("ss.serve.listen"),
		config.GEnvSource("ss_server_rest_listen"),
		config.GConfigSource("server.rest.listen"),
		config.DefaultSource("localhost:8080"),
	)
}
