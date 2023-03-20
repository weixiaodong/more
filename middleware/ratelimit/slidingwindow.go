package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const slidingWindowLimiterTryAcquireRedisScript = `
local key = KEYS[1]
local count = ARGV[1]
local windowTime = ARGV[2]
local time = ARGV[3]
local len = redis.call('llen', key)
if tonumber(len) < tonumber(count) then
	redis.call('lpush', key, time)
	return 1
end

local earlyTime = redis.call('lindex', key, tonumber(len) - 1)
if tonumber(time) - tonumber(earlyTime) < tonumber(windowTime) then
	return 0
end

redis.call('rpop', key)
redis.call('lpush', key, time)
return 1
`

type slidingWindowLimiter struct {
	limit  int           // 窗口请求上限
	window int           // 窗口时间大小
	client *redis.Client // Redis客户端
	script *redis.Script // TryAcquire脚本
}

func newSlidingWindowLimiter(client *redis.Client) *slidingWindowLimiter {
	// // redis过期时间精度最大到毫秒，因此窗口必须能被毫秒整除
	// if window%time.Millisecond != 0 {
	// 	return nil, errors.New("the window uint must not be less than millisecond")
	// }

	limit := 2
	window := 1
	return &slidingWindowLimiter{
		limit:  limit,
		window: window,
		client: client,
		script: redis.NewScript(slidingWindowLimiterTryAcquireRedisScript),
	}
}

func (l *slidingWindowLimiter) Limit(ctx context.Context, resource string) bool {
	now := time.Now().Unix()

	success, err := l.script.Run(ctx, l.client, []string{resource}, l.limit, l.window, now).Bool()
	if err != nil {
		fmt.Println("redis script run error: ", err)
		return false
	}
	// 若到达窗口请求上限，请求失败
	if !success {
		return true
	}
	return false
}
