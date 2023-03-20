package redis

import (
	"context"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/weixiaodong/more/common/log"
)

var (
	redisErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "redis_error_count", Help: "redis error count"},
		[]string{"addr"})
)

func init() {
	prometheus.MustRegister(redisErrorCount)
}

type Client struct {
	*goredis.Client

	opts *goredis.Options
}

func (c *Client) GetGoRedis() *goredis.Client {
	return c.Client
}

func (c *Client) handleCmdErr(ctx context.Context, cmd goredis.Cmder) {
	if err := cmd.Err(); err != nil {
		if err == goredis.Nil {
			return
		}

		log.Error(ctx, "redis_error", "args", cmd.Args(), "err", err)
		redisErrorCount.WithLabelValues(c.opts.Addr).Inc()
	}
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	// 使用 trace client
	cmd := c.Client.Get(ctx, key)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) Set(ctx context.Context, key string, iface interface{}, i time.Duration) (string, error) {
	cmd := c.Client.Set(ctx, key, iface, i)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) HMGet(ctx context.Context, key string, vals ...string) ([]interface{}, error) {
	cmd := c.Client.HMGet(ctx, key, vals...)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) HMSet(ctx context.Context, key string, mdata map[string]interface{}) (bool, error) {
	cmd := c.Client.HMSet(ctx, key, mdata)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) Del(ctx context.Context, vals ...string) (int64, error) {
	cmd := c.Client.Del(ctx, vals...)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) Exists(ctx context.Context, vals ...string) (int64, error) {
	cmd := c.Client.Exists(ctx, vals...)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) Expire(ctx context.Context, key string, i time.Duration) (bool, error) {
	cmd := c.Client.Expire(ctx, key, i)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) SIsMember(ctx context.Context, key string, iface interface{}) (bool, error) {
	cmd := c.Client.SIsMember(ctx, key, iface)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) SAdd(ctx context.Context, key string, vals ...interface{}) (int64, error) {
	cmd := c.Client.SAdd(ctx, key, vals...)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) SDiff(ctx context.Context, vals ...string) ([]string, error) {
	cmd := c.Client.SDiff(ctx, vals...)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) SetBit(ctx context.Context, key string, i int64, i1 int) (int64, error) {
	cmd := c.Client.SetBit(ctx, key, i, i1)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) GetBit(ctx context.Context, key string, i int64) (int64, error) {
	cmd := c.Client.GetBit(ctx, key, i)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) LPush(ctx context.Context, key string, vals ...interface{}) (int64, error) {
	cmd := c.Client.LPush(ctx, key, vals...)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) RPop(ctx context.Context, key string) (string, error) {
	cmd := c.Client.RPop(ctx, key)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	cmd := c.Client.Incr(ctx, key)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) HDel(ctx context.Context, key string, vals ...string) (int64, error) {
	cmd := c.Client.HDel(ctx, key, vals...)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) HExists(ctx context.Context, key string, val string) (bool, error) {
	cmd := c.Client.HExists(ctx, key, val)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) HGet(ctx context.Context, key string, val string) (string, error) {
	cmd := c.Client.HGet(ctx, key, val)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) HSet(ctx context.Context, key string, val string, iface interface{}) (int64, error) {
	cmd := c.Client.HSet(ctx, key, val, iface)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	cmd := c.Client.TTL(ctx, key)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	cmd := c.Client.HGetAll(ctx, key)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) SRandMember(ctx context.Context, key string) (string, error) {
	cmd := c.Client.SRandMember(ctx, key)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) SRandMemberN(ctx context.Context, key string, i int64) ([]string, error) {
	cmd := c.Client.SRandMemberN(ctx, key, i)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) MGet(ctx context.Context, vals ...string) ([]interface{}, error) {
	cmd := c.Client.MGet(ctx, vals...)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) MSet(ctx context.Context, vals ...interface{}) (string, error) {
	cmd := c.Client.MSet(ctx, vals...)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) BitCount(ctx context.Context, key string, pos ...int64) (int64, error) {
	var bc *goredis.BitCount
	if len(pos) == 2 {
		bc = &goredis.BitCount{
			Start: pos[0],
			End:   pos[1],
		}
	}
	cmd := c.Client.BitCount(ctx, key, bc)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) Decr(ctx context.Context, key string) (int64, error) {
	cmd := c.Client.Decr(ctx, key)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) HIncrBy(ctx context.Context, key string, val string, i int64) (int64, error) {
	cmd := c.Client.HIncrBy(ctx, key, val, i)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) HLen(ctx context.Context, key string) (int64, error) {
	cmd := c.Client.HLen(ctx, key)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) LLen(ctx context.Context, key string) (int64, error) {
	cmd := c.Client.LLen(ctx, key)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) LPop(ctx context.Context, key string) (string, error) {
	cmd := c.Client.LPop(ctx, key)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) RPush(ctx context.Context, key string, vals ...interface{}) (int64, error) {
	cmd := c.Client.RPush(ctx, key, vals...)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) SetNX(ctx context.Context, key string, iface interface{}, i time.Duration) (bool, error) {
	cmd := c.Client.SetNX(ctx, key, iface, i)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) SMembers(ctx context.Context, key string) ([]string, error) {
	cmd := c.Client.SMembers(ctx, key)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) SRem(ctx context.Context, key string, vals ...interface{}) (int64, error) {
	cmd := c.Client.SRem(ctx, key, vals...)
	c.handleCmdErr(ctx, cmd)

	return cmd.Result()
}

func (c *Client) Pipeline(ctx context.Context) (iface goredis.Pipeliner) {
	return c.Client.Pipeline()
}
