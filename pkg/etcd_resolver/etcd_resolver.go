package etcd_resolver

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
)

// etcdResolver 实现 gRPC 的 Resolver 接口
type etcdResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	client     *clientv3.Client
	watcher    clientv3.Watcher
	cancelFunc context.CancelFunc
}

// etcdBuilder 实现 gRPC 的 ResolverBuilder 接口
type etcdBuilder struct{}

// 注册 etcd 解析器
func init() {
	//resolver.Register(&etcdBuilder{})
}

// Build 实现 ResolverBuilder 接口
func (b *etcdBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	log.Println("创建 etcd 解析器，目标:", target.URL.String())

	// 创建 etcd 客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{target.URL.Host}, // etcd 地址
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	r := &etcdResolver{
		target: target,
		cc:     cc,
		client: cli,
	}

	// 初始解析
	r.ResolveNow(resolver.ResolveNowOptions{})

	// 启动监听
	r.watch()

	return r, nil
}

// Scheme 实现 ResolverBuilder 接口
func (b *etcdBuilder) Scheme() string {
	return "etcd" // 对应 URI 中的 scheme
}

// ResolveNow 实现 Resolver 接口
func (r *etcdResolver) ResolveNow(options resolver.ResolveNowOptions) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取服务地址列表
	prefix := "/microservices/" + r.target.Endpoint() // 服务注册的前缀
	resp, err := r.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		log.Printf("获取服务地址失败: %v", err)
		return
	}

	//log.Printf("获取服务地址成功: %d 个地址", len(resp.Kvs))

	// 解析地址并更新连接状态
	var addrs []resolver.Address
	for _, kv := range resp.Kvs {
		if len(kv.Value) == 0 {
			log.Printf("发现空地址，跳过: %s", kv.Key)
			continue
		}

		var ins registry.ServiceInstance

		_ = json.Unmarshal(kv.Value, &ins)

		//log.Println("解析到服务实例:", ins)

		for _, endpoint := range ins.Endpoints {
			addrs = append(addrs, resolver.Address{Addr: endpoint})
		}
	}

	_ = r.cc.UpdateState(resolver.State{Addresses: addrs})

	//log.Printf("[%v]更新服务地址: %v", prefix, addrs)
}

// Close 实现 Resolver 接口
func (r *etcdResolver) Close() {
	if r.cancelFunc != nil {
		r.cancelFunc()
	}
	_ = r.client.Close()
}

// 监听服务变化
func (r *etcdResolver) watch() {
	prefix := "/microservices/" + r.target.Endpoint()
	ctx, cancel := context.WithCancel(context.Background())
	r.cancelFunc = cancel

	go func() {
		watchChan := r.client.Watch(ctx, prefix, clientv3.WithPrefix())
		for resp := range watchChan {
			if resp.Err() != nil {
				log.Printf("监听服务变化失败: %v", resp.Err())
				continue
			}

			// 服务地址变更，重新解析
			r.ResolveNow(resolver.ResolveNowOptions{})
		}
	}()
}
