package config

import (
	"github.com/spf13/viper"
)

func IsDev() bool {
	env := viper.GetString("global.env")
	return env == "dev"
}

// 获取服务发现的地址
func GetDiscoveryEndpoints() []string {
	return viper.GetStringSlice("discovery.endpoints")
}

//  获取服务发现名称
func GetDiscoveryServiceNamePrefix() string {
	return viper.GetString("discovery.serviceNamePrefix")
}

// 获取连接服务发现集群的超时时间 单位伟秒
func GetDiscoveryTimeout() int64 {
	return viper.GetInt64("discovery.timeout")
}

// 获取服务注册地址
func GeGrpcServiceAddr() string {
	return viper.GetString("grpcServer.addr")
}

func GetRabbitMQURL() string {
	return viper.GetString("rabitmq.url")
}
