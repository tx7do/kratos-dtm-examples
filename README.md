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
 Target: 'etcd://127.0.0.1:2379/dtmservice' # register dtm server to this url
 EndPoint: 'grpc://localhost:36790'
```

### 配置为Consul服务发现

```yaml
#  dtm: conf.yml
MicroService:
 Driver: 'dtm-driver-kratos' # name of the driver to handle register/discover
 Target: 'consul://127.0.0.1:8500/dtmservice' # register dtm server to this url
 EndPoint: 'grpc://localhost:36790'
```

### 启动DTM服务

```shell
dtm -c ./conf.yml
```
