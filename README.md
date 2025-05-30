# kratos-dtm-examples

## DTM

### Docker部署

```shell
docker run -itd --name dtm -p 36789:36789 -p 36790:36790 yedf/dtm:latest
```

## 二进制安装

### brew

```shell
brew install dtm
```

### go install 安装

```shell
go install github.com/dtm-labs/dtm@latest
```

## 运行DTM服务

### 配置为Etcd服务发现

```yaml
MicroService:
 Driver: 'dtm-driver-kratos' # name of the driver to handle register/discover
 Target: 'etcd://127.0.0.1:2379/dtm-service' # register dtm server to this url
 EndPoint: 'grpc://127.0.0.1:36790'
```

### 配置为Consul服务发现

```yaml
#  dtm: conf.yml
MicroService:
 Driver: 'dtm-driver-kratos' # name of the driver to handle register/discover
 Target: 'consul://127.0.0.1:8500/dtm-service' # register dtm server to this url
 EndPoint: 'grpc://127.0.0.1:36790'
```

### 启动DTM服务

```shell
dtm -c ./conf.yml
```

## 部署Etcd

Etcd Server

```shell
docker run -itd \
    --name etcd-standalone \
    -p 2379:2379 \
    -p 2380:2380 \
    -e ETCDCTL_API=3 \
    -e ALLOW_NONE_AUTHENTICATION=yes \
    -e ETCD_ADVERTISE_CLIENT_URLS=http://0.0.0.0:2379 \
    bitnami/etcd:latest
```

Etcd Keeper

```shell
docker run -d \
  --name etcdkeeper \
  -p 8080:9090 \  # 主机 8080 → 容器 9090
  -e ETCD_KEEPER_PORT=9090 \
  evildecay/etcdkeeper
```

如果为docker部署Etcd Keeper，etcd的地址请填写为：`host.docker.internal:2379`

## 部署Consul

```bash
docker run -itd \
    --name consul-server-standalone \
    -p 8300:8300 \
    -p 8301:8301 \
    -p 8301:8301/udp \
    -p 8500:8500 \
    -p 8600:8600 \
    -p 8600:8600/udp \
    -e CONSUL_BIND_INTERFACE='eth0' \
    -e CONSUL_DISABLE_KEYRING_FILE=true \
    -e CONSUL_AGENT_MODE=server \
    -e CONSUL_ENABLE_UI=true \
    -e CONSUL_BOOTSTRAP_EXPECT=1 \
    -e CONSUL_CLIENT_LAN_ADDRESS=0.0.0.0 \
    bitnami/consul:latest
```

- 管理后台: <http://localhost:8500>
