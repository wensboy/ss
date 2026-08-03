package server

import (
	"context"
	"net"

	"github.com/wensboy/ss/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RpcServer struct {
	lis      net.Listener
	server   *grpc.Server
	Scontext *ServerContext
}

type RpcServerOption func(*RpcServer)
type RpcServerHook = RpcServerOption

func NewRpcServer() *RpcServer {
	return &RpcServer{}
}

func (s *RpcServer) Start() {
	if err := s.server.Serve(s.lis); err != nil {
		panic(err)
	}
}

func (s *RpcServer) Stop() {
	s.server.GracefulStop()
}

func (s *RpcServer) GoStart(ctx context.Context) {
	go func() {
		s.Start()
	}()
	go func() {
		<-ctx.Done()
		s.server.GracefulStop()
	}()
}

func (s *RpcServer) SetListener(lis net.Listener) *RpcServer {
	s.lis = lis
	return s
}

func (s *RpcServer) SetServer(server *grpc.Server) *RpcServer {
	s.server = server
	return s
}

func (s *RpcServer) SetServerContext(scontext *ServerContext) *RpcServer {
	s.Scontext = scontext
	return s
}

func (s *RpcServer) MountModules(opts ...RpcServerHook) {
	for _, opt := range opts {
		opt(s)
	}
}

func (s *RpcServer) MountServers(registerFns ...func(*grpc.Server, *ServerContext)) RpcServerOption {
	if s.server == nil {
		s.server = grpc.NewServer()
	}
	return func(s *RpcServer) {
		for _, fn := range registerFns {
			fn(s.server, s.Scontext)
		}
	}
}

// 执行在挂载 server 之前, 主要用于挂载拦截器
func (s *RpcServer) MountInterceptors(unaryInterceptors []grpc.UnaryServerInterceptor, streamInterceptors []grpc.StreamServerInterceptor) RpcServerOption {
	return func(s *RpcServer) {
		s.server = grpc.NewServer(
			grpc.ChainUnaryInterceptor(unaryInterceptors...),
			grpc.ChainStreamInterceptor(streamInterceptors...),
		)
	}
}

func (s *RpcServer) MountDatabases(preHook, runHook, postHook RpcServerHook) RpcServerOption {
	return func(s *RpcServer) {
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

func (s *RpcServer) MountConfig(preHook, runHook, postHook RpcServerHook) RpcServerOption {
	return func(s *RpcServer) {
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

func (s *RpcServer) mountConfig() {
	enableReflection := config.MustLookup[bool](
		config.GEnvSource("ss_server_rpc_reflection"),
		config.GConfigSource("server.rpc.reflection"),
		config.DefaultSource(false),
	)
	if enableReflection {
		// 用于 grpcurl 等工具调试
		reflection.Register(s.server)
	}
	lisNetwork := config.MustLookup[string](
		config.GEnvSource("ss_server_rpc_network"),
		config.GConfigSource("server.rpc.network"),
		config.DefaultSource("tcp"),
	)
	lisAddr := config.MustLookup[string](
		config.GFlagSource("ss.serve.listen"),
		config.GEnvSource("ss_server_rpc_listen"),
		config.GConfigSource("server.rpc.listen"),
		config.DefaultSource("localhost:50051"),
	)
	var err error
	s.lis, err = net.Listen(lisNetwork, lisAddr)
	if err != nil {
		panic(err)
	}
}
