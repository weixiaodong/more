package service

import (
    "context"
    "log"

    "github.com/weixiaodong/more/protos/pb"
)

const (
    port = ":50051"
)

// HelloService is used to implement helloworld.GreeterServer.
type HelloService struct{}

// SayHello implements helloworld.GreeterServer
func (s *HelloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
    log.Printf("Received: %v", in.Name)
    return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
