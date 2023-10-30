package grpcserver

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"gitee.com/git-lz/go-tinyid/common/xgrpc/proto"
)

type IdSequenceServer struct {
}

func NewIdSequenceServer() *IdSequenceServer {
	return &IdSequenceServer{}
}

func Init() func() {
	listen, err := net.Listen("tcp", "127.0.0.1:8098")
	if err != nil {
		panic(err)
	}

	server := New(WithMaxConcurrentStreams(100))
	proto.RegisterIdSequenceServiceServer(server, NewIdSequenceServer())

	go func() {
		err := server.Serve(listen)
		if err != nil {
			panic(err)
		}
	}()
	fmt.Println("init proxy grpc server success!")

	return func() {
		server.GracefulStop()
		fmt.Println("stop grpc server success...")
	}
}

func New(opts ...Option) *grpc.Server {
	o := newDefaultServerOptions()

	for _, opt := range opts {
		opt(o)
	}

	s := grpc.NewServer(
		grpc.MaxConcurrentStreams(o.maxConcurrentStreams),
	)

	reflection.Register(s)

	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	return s
}
