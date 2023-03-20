package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
)

var (
	ErrLockFailed = errors.New("get lock failed with max retry times")
	ErrTimeout    = errors.New("get lock failed with timeout")
)

type redisLock struct {
	resource        string
	randomValue     string
	ttl             time.Duration
	tryLockInterval time.Duration
	watchDog        chan struct{}
}

type OptionFunc func(*redisLock)

func WithWatchDog() OptionFunc {
	return func(l *redisLock) {
		l.watchDog = make(chan struct{})
	}
}

func WithTryLockInterval(t time.Duration) OptionFunc {
	return func(l *redisLock) {
		l.tryLockInterval = t
	}
}

func NewLock(key string, duration time.Duration, opts ...OptionFunc) *redisLock {
	l := &redisLock{
		resource: key,
		ttl:      duration,
	}
	for _, opt := range opts {
		opt(l)
	}
	return l

}

func (l *redisLock) TryLock(ctx context.Context) error {
	randomValue := xid.New().String()
	success, err := GetClient().SetNX(ctx, l.resource, randomValue, l.ttl)
	if err != nil {
		return err
	}
	// 加锁失败
	if !success {
		return ErrLockFailed
	}
	// 加锁成功
	l.randomValue = randomValue

	if l.watchDog != nil {
		go l.startWatchDog()
	}
	return nil
}

var unlockScript = redis.NewScript(`
if redis.call("get",KEYS[1]) == ARGV[1] then
    return redis.call("del",KEYS[1])
else
    return 0
end
`)

func (l *redisLock) Unlock(ctx context.Context) error {
	fmt.Println("Unlock")
	if l.watchDog != nil {
		close(l.watchDog)
	}
	return unlockScript.Run(context.Background(), GetClient().GetGoRedis(), []string{l.resource}, l.randomValue).Err()
}

func (l *redisLock) Lock(ctx context.Context) error {
	// 尝试加锁
	err := l.TryLock(ctx)
	if err == nil {
		return nil
	}
	if !errors.Is(err, ErrLockFailed) {
		return err
	}

	// 加锁失败，不断尝试
	if l.tryLockInterval == 0 {
		return err
	}

	ticker := time.NewTicker(l.tryLockInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			// 超时
			return ErrTimeout
		case <-ticker.C:
			// 重新尝试加锁
			err := l.TryLock(ctx)
			if err == nil {
				return nil
			}
			if !errors.Is(err, ErrLockFailed) {
				return err
			}
		}
	}
}

func (l *redisLock) startWatchDog() {
	ticker := time.NewTicker(l.ttl / 3)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// 延长锁的过期时间
			ctx, cancel := context.WithTimeout(context.Background(), l.ttl/3*2)
			fmt.Println("expire")
			ok, err := GetClient().Expire(ctx, l.resource, l.ttl)
			cancel()
			// 异常或锁已经不存在则不再续期
			if err != nil || !ok {
				return
			}
		case <-l.watchDog:
			fmt.Println("done")

			// 已经解锁
			return
		}
	}
}
