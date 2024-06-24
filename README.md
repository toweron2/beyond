```shell
#!/usr/bin/env bash
/d/environment/etcd-v3.5.10-windows-amd64/etcd.exe &
/d/environment/redis/redis-server.exe &

```


# beyond

## 微信

15951703783 备注: 我爱golang

## 架构图

### 第一课 (项目概述)
#### [文档](https://pwmzlkcu3p.feishu.cn/docx/HsMIdfEa0ogEvmxCpRrcEOyenqc) [视频](https://www.bilibili.com/video/BV1op4y177iS/)


### 第二课  (微服务拆分&&项目结构 && 服务初始化 && 调用流程 && jwt验证 && 验证码注册 && 缓存 && 服务注册与发现)
#### [文档](https://pwmzlkcu3p.feishu.cn/docx/XX6xdpB0UoH0auxgPYlcxmninDb)   [视频](https://www.bilibili.com/video/BV1CH4y1Q7PM/)

```shell
#查看etcd中是否已注册，查看go-zero源码，看看在哪里注册的，怎么注册的
etcdctl get --prefix user.rpc

# 查看租约
etcdctl lease list

# 查看租约剩余时间
etcdctl lease timetolive 694d8a5192174527
```

### 第三课 (自定义业务错误码，RPC和API服务通用错误码)
#### [文档](https://pwmzlkcu3p.feishu.cn/docx/XLB9dK9Cao3Z7HxPyEscu0P9nIb) [视频](https://www.bilibili.com/video/BV19u411w7WS/)

### 第四课 实现文章功能和互动功能(架构设计&表设计)
#### [文档](https://pwmzlkcu3p.feishu.cn/docx/U9FGdVAFuoFFiUxySsgcl5TMnke) [视频](https://www.bilibili.com/video/BV1Y8411q7uW/)

### 第五课
#### [文档](https://pwmzlkcu3p.feishu.cn/docx/EBzWdSFR5oPVMJxP1oOcUGojnTd) [视频](https://www.bilibili.com/video/BV1k8411y7W5/)


### 第六课
#### [文档](https://pwmzlkcu3p.feishu.cn/docx/O1u3d9pqWo4sZTx6NTxcOHhknmd) [视频](https://www.bilibili.com/video/BV1F84y1S74g/)
```shell
# 创建网络
docker network create beyond --driver bridge
# 拉取zookeeper
docker pull zookeeper
# 创建zookeeper容器
docker run -d --name zookeeper --network beyond -p 2181:2181 -t zookeeper
# 拉取kafka
docker pull wurstmeister/kafka
# 创建kafka容器
docker run -d --name kafka --network beyond -p 9092:9092 -e KAFKA_BROKER_ID=0 \
  -e KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181 \
  -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://192.168.5.113:9092 \
  -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092 wurstmeister/kafka
  
  
192.168.5.113
# 进入kafka容器
docker exec -it {{容器ID}} /bin/bash
# 进入kafka执行命令目录
cd /opt/kafka/bin
# 创建topic
./kafka-topics.sh --create --topic topic-beyond-like --bootstrap-server localhost:9092
./kafka-topics.sh --create --topic topic-beyond-like --bootstrap-server 192.168.5.113:9092
# 查看topic信息
./kafka-topics.sh --describe --topic topic-beyond-like --bootstrap-server localhost:9092
./kafka-topics.sh --describe --topic topic-beyond-like --bootstrap-server 192.168.5.113:9092
# 生产消息
./kafka-console-producer.sh --topic topic-beyond-like --bootstrap-server localhost:9092
./kafka-console-producer.sh --topic topic-beyond-like --bootstrap-server 192.168.5.113:9092
# 消费消息
./kafka-console-consumer.sh --topic topic-beyond-like --from-beginning \
--bootstrap-server localhost:9092
--bootstrap-server 92.168.5.113:9092
# 查看所有topic
kafka-topics.sh --list --bootstrap-server localhost:9092

```

###### 生产者
```yaml
# 配置
KqPusherConf:
  Brokers:
    - 127.0.0.1:9092
  Topic: topic-beyond-like
```

```go
// 初始化
type ServiceContext struct {
    Config         config.Config
    KqPusherClient *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config:         c,
        KqPusherClient: kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
    }
}
```

```go
// 发送kafka消息
func (l *ThumbupLogic) Thumbup(in *service.ThumbupRequest) (*service.ThumbupResponse, error) {
    // TODO 逻辑暂时忽略
    // 1. 查询是否点过赞
    // 2. 计算当前内容的总点赞数和点踩数
    
    msg := &types.ThumbupMsg{
        BizId:    in.BizId,
        ObjId:    in.ObjId,
        UserId:   in.UserId,
        LikeType: in.LikeType,
    }
    // 发送kafka消息，异步
    threading.GoSafe(func() {
        data, err := json.Marshal(msg)
        if err != nil {
            l.Logger.Errorf("[Thumbup] marshal msg: %+v error: %v", msg, err)
            return
        }
        err = l.svcCtx.KqPusherClient.Push(string(data))
        if err != nil {
            l.Logger.Errorf("[Thumbup] kq push data: %s error: %v", data, err)
        }
    })
    
    return &service.ThumbupResponse{}, nil
}

```
###### 消费者
```yaml
# 配置
Name: mq
KqConsumerConf:
  Name: like-kq-consumer
  Brokers:
    - 127.0.1:9092
  Group: group-beyond-like
  Topic: topic-beyond-like
  Offset: last
  Consumers: 1
  Processors: 1
```
```go
// 消费kafka信息
type ThumbupLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
    return &ThumbupLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

func (l *ThumbupLogic) Consume(key, val string) error {
    fmt.Printf("get key: %s val: %s\n", key, val)
    return nil
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
    return []service.Service{
        kq.MustNewQueue(svcCtx.Config.KqConsumerConf, NewThumbupLogic(ctx, svcCtx)),
    }
}
```

###### grpcurl工具
```go
// 启动反射服务
func main() {
   flag.Parse()

   var c config.Config
   conf.MustLoad(*configFile, &c)
   ctx := svc.NewServiceContext(c)

   s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
      service.RegisterLikeServer(grpcServer, server.NewLikeServer(ctx))

      if c.Mode == zs.DevMode || c.Mode == zs.TestMode {
         reflection.Register(grpcServer)
      }
   })
   defer s.Stop()

   fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
   s.Start()
}
```
```yaml
# 配置
Mode: test
```
```shell
# 安装
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
# 查看服务列表
grpcurl -plaintext 127.0.0.1:8080 list
# 查看服务方法列表
grpcurl -plaintext 127.0.0.1:8080 list service.Like
# 调用方法
grpcurl -plaintext -d '{"bizId": "article", "objId": 123, "userId": 234, "likeType": 1}' \
 127.0.0.1:8080 service.Like/Thumbup
```



### 第七课
#### [文档](https://pwmzlkcu3p.feishu.cn/docx/SvSwdAETNo3F0Axznp8ceaX1nIb) [视频](https://www.bilibili.com/video/BV1sz4y1G73u/)

### 第八课
#### [文档](https://pwmzlkcu3p.feishu.cn/docx/At4zdJzrFowGMmx5OjJcrS5mnir) [视频](https://www.bilibili.com/video/BV1S8411C7CY/)

### 第九课
#### [文档](https://pwmzlkcu3p.feishu.cn/docx/OW6Sd0ZsioU9LLxuUSXcUZgNnoD) [视频](https://www.bilibili.com/video/BV1oB4y1f7Tr/)

### 第十课
#### [文档](https://pwmzlkcu3p.feishu.cn/docx/Si1Cd4EGxoZXkJxGenzcFttOnsh) [视频](https://www.bilibili.com/video/BV1je411R7iy/)

### 第十一课
#### [文档](https://pwmzlkcu3p.feishu.cn/docx/NYyHdpSzhoB8Zkxdb6NcpUGynH4) [视频](https://www.bilibili.com/video/BV11u4y1Y7GC/)

### 第十二课
#### [文档](https://pwmzlkcu3p.feishu.cn/docx/QZcKdB4VXoUCRDxGilfcTaw7n8e) [视频](https://www.bilibili.com/video/BV1u64y177rL)

### 第十三课
#### [文档](https://pwmzlkcu3p.feishu.cn/docx/BDNUdhmP4oec6ix1P1ZcqF3rnIc) [视频](https://www.bilibili.com/video/BV1bH4y1C7Uj)

### 第十四课
#### [文档](https://pwmzlkcu3p.feishu.cn/docx/Ydd4dG8OSobJ1rxJFgacOBKvnSg)