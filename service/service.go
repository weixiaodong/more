package service

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/weixiaodong/more/common/redis"
    "github.com/weixiaodong/more/protos/pb"
)

// HelloService is used to implement helloworld.GreeterServer.
type HelloService struct{}

// SayHello implements helloworld.GreeterServer
func (s *HelloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {

    l := redis.NewLock(fmt.Sprintf("lock:%s", in.Name), 30*time.Second, redis.WithWatchDog())
    // 如果请求context未设置过期时间，设置锁过期时间
    ctx, _ = context.WithTimeout(ctx, 5*time.Second)
    err := l.Lock(ctx)
    if err != nil {
        return nil, err
    }
    defer l.Unlock(ctx)

    time.Sleep(60 * time.Second)
    log.Printf("Received: %v", in.Name)
    return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
