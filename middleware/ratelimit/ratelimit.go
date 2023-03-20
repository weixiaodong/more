package ratelimit

import (
	"context"
	"fmt"
	"path"

	"github.com/go-redis/redis_rate/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/weixiaodong/more/common/redis"
)

// Limiter defines the interface to perform request rate limiting.
// If Limit function return true, the request will be rejected.
// Otherwise, the request will pass.
type Limiter interface {
	Limit(ctx context.Context, resource string) bool
}

type goRedisLimiter struct {
	*redis_rate.Limiter
}

func (l *goRedisLimiter) Limit(ctx context.Context, method string) bool {
	limitKey := path.Base(method)
	fmt.Println("limitKey: ", limitKey)
	res, err := l.Allow(context.Background(), limitKey, redis_rate.PerSecond(2))
	if err != nil {
		panic(err)
	}
	fmt.Println("allowed", res.Allowed, "remaining", res.Remaining)

	if res.Allowed > 0 || res.Remaining > 0 {
		return false
	}
	return true
}

type methodLimiter struct {
	Limiter
}

func (l *methodLimiter) Limit(ctx context.Context, method string) bool {
	return l.Limit(ctx, method)
}

func NewLimiter() Limiter {
	goRedisClient := redis.GetClient().GetGoRedis()
	// return &goRedisLimiter{
	// 	Limiter: redis_rate.NewLimiter(goRedisClient),
	// }
	// return newFixedWindowLimiter(goRedisClient)
	return newSlidingWindowLimiter(goRedisClient)

}

// UnaryServerInterceptor returns a new unary server interceptors that performs request rate limiting.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	limiter := NewLimiter()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		fmt.Println("limit start")
		method := info.FullMethod
		if limiter.Limit(ctx, method) {
			return nil, status.Errorf(codes.ResourceExhausted, "%s is rejected by grpc_ratelimit middleware, please retry later.", info.FullMethod)
		}
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new stream server interceptor that performs rate limiting on the request.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	limiter := NewLimiter()
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		method := info.FullMethod
		if limiter.Limit(context.Background(), method) {
			return status.Errorf(codes.ResourceExhausted, "%s is rejected by grpc_ratelimit middleware, please retry later.", info.FullMethod)
		}
		return handler(srv, stream)
	}
}
