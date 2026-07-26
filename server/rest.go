package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/wensboy/ss/config"
	"github.com/wensboy/ss/db"
)

const (
	GLOBAL_MIDDLEWARE = "_global_middleware"
)

type RestServer struct {
	muxer    *echo.Echo
	server   *http.Server
	routers  map[string]*echo.Group // 用于中间件追踪挂载
	scontext *ServerContext
}

func NewRestServer() *RestServer {
	return &RestServer{
		muxer:    echo.New(),
		server:   &http.Server{},
		routers:  make(map[string]*echo.Group),
		scontext: NewServerContext(),
	}
}

type RestServerOption func(*RestServer)
type RestServerHook = RestServerOption

func (s *RestServer) Start() {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func (s *RestServer) Stop() {
	if err := s.server.Shutdown(context.Background()); err != nil {
		panic(err)
	}
}

func (s *RestServer) GoStart(ctx context.Context) {
	go func() {
		s.Start()
	}()
	go func() {
		<-ctx.Done()
		if err := s.server.Shutdown(context.Background()); err != nil {
			panic(err)
		}
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

func (s *RestServer) MountModules(opts ...RestServerHook) {
	for _, opt := range opts {
		opt(s)
	}
}

func (s *RestServer) MountRouters(prefix string, registers ...func(*echo.Group)) RestServerOption {
	if s.routers[prefix] == nil {
		s.routers[prefix] = s.muxer.Group(prefix)
	}
	return func(s *RestServer) {
		for _, register := range registers {
			register(s.routers[prefix])
		}
	}
}

// 总是在 router 之前挂载中间件
func (s *RestServer) MountMiddlewares(prefix string, middlewares ...echo.MiddlewareFunc) RestServerOption {
	if s.routers[prefix] == nil {
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
			s.mountDatabases()
		}
		if postHook != nil {
			postHook(s)
		}
	}
}

func (s *RestServer) mountDatabases() {
	// sql database
	{
		dbType := config.MustLookup[string](
			config.GEnvSource("ss_db_type"),
			config.GConfigSource("db.type"),
			config.DefaultSource(db.DB_TYPE_SQLITE),
		)
		dbname := config.MustLookup[string](
			config.GEnvSource("ss_db_name"),
			config.GConfigSource("db.name"),
			config.DefaultSource("default"),
		)
		dsn := config.MustLookup[string](
			config.GEnvSource("ss_db_dsn"),
			config.GConfigSource("db.dsn"),
		)
		sqlDB, err := db.NewSqlDatabase(dbType, dbname, dsn)
		if err != nil {
			panic(err)
		}
		db.GetGSqlDBContext().Set(sqlDB)
		s.scontext.SetDBContext(db.GetGSqlDBContext())
	}
}

func (s *RestServer) MountConfig(preHook, runHook, postHook RestServerHook) RestServerOption {
	if s.server == nil || s.muxer == nil {
		// todo: 创建带编码错误 Err
		panic("[server] server or muxer is nil")
	}
	s.server.Handler = s.muxer
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
