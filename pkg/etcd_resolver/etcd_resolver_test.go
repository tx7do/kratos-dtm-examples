package etcd_resolver

import (
	"context"
	"testing"

	"github.com/dtm-labs/client/dtmgrpc/dtmgpb"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc/resolver/discovery"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Test_CallGid_EtcdBuilder(t *testing.T) {
	// EndPoint: '127.0.0.1:36790'
	// 不能够为：EndPoint: 'grpc://127.0.0.1:36790'

	resolver.Register(&etcdBuilder{})

	// 直接连接 etcd 服务发现的 DTM
	conn, err := grpc.NewClient(
		"etcd://127.0.0.1:2379/dtm-service",
		//"discovery:///dtm-service",
		//"127.0.0.1:36790",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), // 负载均衡策略
	)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	if conn == nil {
		log.Fatal("连接失败: 连接对象为 nil")
		return
	}

	client := dtmgpb.NewDtmClient(conn)
	// 调用 GenGid 方法
	out, err := client.NewGid(context.Background(), &emptypb.Empty{})
	if err != nil {
		t.Fatalf("调用 NewGid 方法失败: %v", err)
	} else {
		t.Logf("调用 NewGid 方法成功, GID: %s", out.Gid)
	}
}

func Test_CallGid_KratosBuilder(t *testing.T) {
	d := bootstrap.NewDiscovery(
		&conf.Registry{
			Type: "etcd",
			Etcd: &conf.Registry_Etcd{
				Endpoints: []string{"127.0.0.1:2379"},
			},
		},
	)
	resolver.Register(discovery.NewBuilder(d, discovery.WithInsecure(true)))

	//resolver.Register(&etcdBuilder{})

	// 直接连接 etcd 服务发现的 DTM
	conn, err := grpc.NewClient(
		//"etcd://127.0.0.1:2379/dtm-service",
		"discovery:///dtm-service",
		//"127.0.0.1:36790",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), // 负载均衡策略
	)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	if conn == nil {
		log.Fatal("连接失败: 连接对象为 nil")
		return
	}

	client := dtmgpb.NewDtmClient(conn)
	if client == nil {
		log.Fatal("连接失败: 客户端对象为 nil")
		return
	}

	// 调用 GenGid 方法
	out, err := client.NewGid(context.Background(), &emptypb.Empty{})
	if err != nil {
		t.Fatalf("调用 NewGid 方法失败: %v", err)
	} else {
		t.Logf("调用 NewGid 方法成功, GID: %s", out.Gid)
	}
}
