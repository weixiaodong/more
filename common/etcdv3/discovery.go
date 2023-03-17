package etcdv3

import (
	"context"
	"log"
	"sync"
	"time"

	"go.etcd.io/etcd/client/v3"
)

// ServiceDiscovery 服务发现
type ServiceDiscovery struct {
	cli        *clientv3.Client
	serverList map[string]string // 服务列表
	lock       sync.Mutex
}

// NewServiceDiscovery 发现服务
func NewServiceDiscovery(endpoint []string) *ServiceDiscovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoint,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	return &ServiceDiscovery{
		cli:        cli,
		serverList: make(map[string]string),
	}
}

// WatchService 初始化服务列表与监视
func (s *ServiceDiscovery) WatchService(prefix string) error {
	// 根据前缀获取现有的key
	resp, err := s.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, ev := range resp.Kvs {
		s.SetServerList(string(ev.Key), string(ev.Value))
	}

	// 监听后续操作
	go s.watcher(prefix)

	return nil
}

// SetServerList 新增服务地址
func (s *ServiceDiscovery) SetServerList(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.serverList[key] = val
	log.Println("put key: ", key, "val: ", val)
}

// 监听前缀
func (s *ServiceDiscovery) watcher(prefix string) {
	rch := s.cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	log.Printf("watching prefix: %s now...", prefix)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				// 修改或新增
				s.SetServerList(string(ev.Kv.Key), string(ev.Kv.Value))
			case clientv3.EventTypeDelete:
				// 删除
				s.DelServerList(string(ev.Kv.Key))
			}
		}
	}
}

// DelServerList 删除服务地址
func (s *ServiceDiscovery) DelServerList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.serverList, key)
	log.Println("delete key: ", key)
}

// GetServices 获取服务地址
func (s *ServiceDiscovery) GetServices() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	addrs := make([]string, 0)
	for _, v := range s.serverList {
		addrs = append(addrs, v)
	}
	return addrs
}

func (s *ServiceDiscovery) Close() {
	s.cli.Close()
}
