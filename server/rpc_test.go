package server

import (
	"context"
	"log"
	"testing"

	commonPb "github.com/wensboy/ss/proto/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type MockRpcServer struct {
	commonPb.UnimplementedCommonServiceServer
}

func (s *MockRpcServer) Ping(_ context.Context, req *commonPb.PingReq) (*commonPb.PingResp, error) {
	var msg string = "pon"
	if req.Msg == "ping" {
		msg = "pong"
	}
	return &commonPb.PingResp{
		Msg: msg,
	}, nil
}

func mountMockRpcServer(s *grpc.Server, sctx *ServerContext) {
	commonPb.RegisterCommonServiceServer(s, &MockRpcServer{})
}

func logUnary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Printf("method: %s req: %#v\n", info.FullMethod, req)
	return handler(ctx, req)
}

func Test_RpcServer(t *testing.T) {
	MockLoadAll()
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(logUnary),
	)
	reflection.Register(server)
	rpcServer := NewRpcServer().SetServer(server)
	rpcServer.MountModules(
		rpcServer.MountServers(mountMockRpcServer),
		rpcServer.MountConfig(nil, nil, nil),
	)
	rpcServer.Start()
}
