1. 项目使用 cobra 命令目录组织方式

2. 目录结构与功能

```bash
├── cmd
│ ├── client.go           # 模拟grpc客户端
│ ├── root.go
│ └── serve.go            # 启动grpc服务
├── common                # 服务公共依赖库
│ ├── config              # 配置
│ ├── etcdv3              # etcd客户端封装
│ ├── log                 # log封装
│ ├── rabbitmq            # rabbitmq客户端封装
│ └── redis               # redis客户端封装
├── etc
│ └── config.yaml
├── go.mod
├── go.sum
├── LICENSE
├── main.go               # 服务启动
├── Makefile
├── middleware
│ ├── chain.go            # 服务中间件链工具
│ ├── ratelimit           # 限流中间件
│ └── recovery
├── protos                # 接口协议
│ ├── pb
│ └── service.proto
└── service
    └── service.go        # 服务逻辑
```
