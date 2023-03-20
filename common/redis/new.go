package redis

import (
	"sync"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var (
	once sync.Once
	cli  *Client
)

func GetClient() *Client {

	once.Do(func() {

		cli = &Client{}
		cli.opts = &goredis.Options{
			Addr:         viper.GetString("redis.addr"),
			Password:     viper.GetString("redis.pwd"),
			DB:           viper.GetInt("redis.db"),
			DialTimeout:  time.Duration(viper.GetFloat64("redis.dialTimeout")) * time.Second,
			ReadTimeout:  time.Duration(viper.GetFloat64("redis.readTimeout")) * time.Second,
			WriteTimeout: time.Duration(viper.GetFloat64("redis.writeTimeout")) * time.Second,
			PoolSize:     viper.GetInt("redis.poolSize"),
			MinIdleConns: viper.GetInt("redis.minIdleConns"),
			MaxRetries:   viper.GetInt("redis.maxRetries"),
		}

		cli.Client = goredis.NewClient(cli.opts)
	})

	return cli
}
