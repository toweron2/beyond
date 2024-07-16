**beyond**

------

### 服务启动

```shell
#!/usr/bin/env bash
# etcd
/d/environment/etcd-v3.5.10-windows-amd64/etcd.exe &
# etcdkeeper
/d/environment/etcdkeeper/etcdkeeper.exe -p 2380 &
# redis
/d/environment/redis/redis-server.exe &
# zookeeper
/d/environment/zookeeper-3.9.2/bin/zkServer.cmd &
# kafka
/d/environment/kafka/bin/windows/kafka-server-start.bat /d/environment/kafka/config/server.properties &
# canal
# /d/environment/canal.deployer-1.1.8/bin/stop.sh
/d/environment/canal.deployer-1.1.8/bin/startup.sh &
# prometheus
/d/environment/prometheus-2.53.0.windows-amd64/prometheus.exe --config.file=/d/environment/prometheus-2.53.0.windows-amd64/prometheus.yml &
# grafana
/d/environment/grafana-v11.1.0/bin/grafana-server.exe --homepath=/d/environment/grafana-v11.1.0 &
# jaeger
/d/environment/jaeger-1.58.0-windows-amd64/jaeger-all-in-one.exe & 
# elasticsearch
/d/environment/elasticsearch-8.14.2/bin/elasticsearch.bat &
# kibana
/d/environment/kibana-8.14.2/bin/kibana.bat &

user.rpc 6001
applet.api 8888
article.rpc 6101
article.api 80
article.mq
like.rpc 6201
like.mq
follow.rpc 6301
```



```tex
微信
15951703783 备注: 我爱golang
```





#### 架构图



### 第一课-项目概述

文档地址: https://pwmzlkcu3p.feishu.cn/docx/HsMIdfEa0ogEvmxCpRrcEOyenqc
视频地址: https://www.bilibili.com/video/BV1op4y177iS/






### 第二课-微服务拆分&&项目结构 && 服务初始化 && 调用流程 && jwt验证 && 验证码注册 && 缓存 && 服务注册与发现

文档地址: https://pwmzlkcu3p.feishu.cn/docx/XX6xdpB0UoH0auxgPYlcxmninDb)   

视频地址: https://www.bilibili.com/video/BV1CH4y1Q7PM/

```shell
#查看etcd中是否已注册，查看go-zero源码，看看在哪里注册的，怎么注册的
etcdctl get --prefix user.rpc

# 查看租约
etcdctl lease list

# 查看租约剩余时间
etcdctl lease timetolive 694d8a5192174527
```



### 第三课-自定义业务错误码，RPC和API服务通用错误码

文档地址: https://pwmzlkcu3p.feishu.cn/docx/XLB9dK9Cao3Z7HxPyEscu0P9nIb
视频地址: https://www.bilibili.com/video/BV19u411w7WS/



### 第四课-实现文章功能和互动功能(架构设计&表设计)

文档地址: https://pwmzlkcu3p.feishu.cn/docx/U9FGdVAFuoFFiUxySsgcl5TMnke
视频地址: https://www.bilibili.com/video/BV1Y8411q7uW/



### 第五课-项目结构说明 && 文件上传 && 对接阿里OSS && 文章发布 && 多服务联动测试

文档地址: https://pwmzlkcu3p.feishu.cn/docx/EBzWdSFR5oPVMJxP1oOcUGojnTd
视频地址: https://www.bilibili.com/video/BV1k8411y7W5/




### 第六课-docker安装Kafka & 在go-zero中使用kafka进行数据生产和消费 & grpcurl工具使用

文档地址: https://pwmzlkcu3p.feishu.cn/docx/O1u3d9pqWo4sZTx6NTxcOHhknmd
视频地址: https://www.bilibili.com/video/BV1F84y1S74g/

##### Kafka安装

zookeeper下载地址: https://dlcdn.apache.org/zookeeper/zookeeper-3.9.2/

kafka下载地址: https://dlcdn.apache.org/kafka/3.7.1/

```sh
# zookeeper配置
zookeeper 根目录下创建 data文件加
# 复制conf 目录下 zoo_sample.cfg文件，改名为zoo.cfg
# 将dataDir=/tmp/zookeeper
dataDir=D:\\environment\\apache-zookeeper-3.9.2-bin\\data

# kafka配置
# kafka根目录创建logs文件夹
# 修改config/server.properties
# 修改 log.dirs 参数值，修改成上一步新建的logs文件夹
log.dirs=D:\\environment\\kafka_2.13-3.7.1\\logs
# 修改 listeners 参数
listeners=PLAINTEXT://localhost:9092


# 启动
zkServer
kafka-server-start.bat .\config\server.properties
```



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



### 第七课-Canal安装配置 & Mysql配置 & Canal解析Binlog投递到Kafka & Kafka消费Mysql数据变更事件

文档地址: https://pwmzlkcu3p.feishu.cn/docx/SvSwdAETNo3F0Axznp8ceaX1nIb
视频地址: https://www.bilibili.com/video/BV1sz4y1G73u/

##### Canel

下载地址:  https://github.com/alibaba/canal/releases

```shell
# 安装
mkdir -p /usr/local/canal

cp canal.deployer-xxx.tar.gz /usr/local/canal

tar -zxvf canal.deployer-xxx.tar.gz 
```

- Mysql配置

  对于自建 MySQL , 需要先开启 Binlog 写入功能，配置 binlog-format 为 ROW 模式，my.cnf 中配置如下：

  ```yaml
  [mysqld]
  log-bin=mysql-bin # 开启 binlog
  binlog-format=ROW # 选择 ROW 模式
  server_id=1 # 配置 MySQL replaction 需要定义，不要和 canal 的 slaveId 重复
  ```

  ```sql
  -- 权 canal 链接 MySQL 账号具有作为 MySQL slave 的权限, 如果已有账户可直接 grant
  CREATE USER canal IDENTIFIED BY 'canal';  
  GRANT SELECT, REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'canal'@'%';
  -- GRANT ALL PRIVILEGES ON *.* TO 'canal'@'%' ;
  FLUSH PRIVILEGES;
  ```

- canal instance配置

  ```shell
  vim /usr/local/canal/conf/example/instance.properties
  ```

  ```properties
  # position info
  canal.instance.master.address=127.0.0.1:3306
  
  # username/password
  canal.instance.dbUsername=canal
  canal.instance.dbPassword=canal
  
  # mq config
  canal.mq.topic=topic-like-count
  canal.mq.partitionsNum=1
  #库名.表名: 唯一主键，多个表之间用逗号分隔
  canal.mq.partitionHash=beyond_like.like_count:id
  ```

- cnaal配置

  ```shell
  vim /usr/local/canal/conf/canal.properties
  ```

  ```properties
  # tcp, kafka, rocketMQ, rabbitMQ, pulsarMQ
  canal.serverMode = kafka
  
  ##################################################
  #########                    Kafka                   #############
  ##################################################
  kafka.bootstrap.servers = 127.0.0.1:9092
  ```
  
  ```shell
  # 启动
  cd /usr/local/canal/
  
  sh bin/startup.sh
  
  # 停止
  cd /usr/local/canal/
  
  sh bin/stop.sh
  ```
  
  


----

### 第八课-文章列表缓存 && 缓存代码 & 性能提升工具MapReduce

文档地址: https://pwmzlkcu3p.feishu.cn/docx/At4zdJzrFowGMmx5OjJcrS5mnir
视频地址: https://www.bilibili.com/video/BV1S8411C7CY/

### 第九课-缓存一致性保证 & 缓存击穿 & 缓存穿透 & 缓存雪崩 

文档地址: https://pwmzlkcu3p.feishu.cn/docx/OW6Sd0ZsioU9LLxuUSXcUZgNnoD
视频地址: https://www.bilibili.com/video/BV1oB4y1f7Tr/

### 第十课-在go-zero中集成GORM & 基于GORM实现关注服务核心功能 & GORM集成指标监控和链路追踪 & 服务对接Prometheus & 服务对接Jaeger

文档地址: https://pwmzlkcu3p.feishu.cn/docx/Si1Cd4EGxoZXkJxGenzcFttOnsh
视频地址: https://www.bilibili.com/video/BV1je411R7iy/



##### 架构图

##### go-zero集成GORM

什么叫集成？

尽量复用go-zero能力，比如日志等；要和go-zero的其他组件模块串联一体，比如Trace；提供和go-zero其他组件模块相同的Metrics能力。总之，就要像是和使用go-zero框架自带的组件一样。

##### 定义GORM结构体

```go
type (
	Config struct {
		DSN          string
		MaxOpenConns int `json:",default=10"`
		MaxIdleConns int `json:",default=100"`
		MaxLifetime  int `json:",default=3600"`
	}
	DB struct {
		*gorm.DB
	}
	ormLog struct {
		LogLevel logger.LogLevel
	}
)
```

###### 使用go-zero日志

```go
func (l *ormLog) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

func (l *ormLog) Info(ctx context.Context, format string, v ...interface{}) {
	if l.LogLevel <= logger.Info {
		return
	}
	logx.WithContext(ctx).Infof(format, v...)
}

func (l *ormLog) Warn(ctx context.Context, format string, v ...interface{}) {
	logx.WithContext(ctx).Infof(format, v...)
}

func (l *ormLog) Error(ctx context.Context, format string, v ...interface{}) {
	logx.WithContext(ctx).Errorf(format, v...)
}

func (l *ormLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	logx.WithContext(ctx).WithDuration(elapsed).Infof("[%.3fms] [rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
}
```

###### 创建DB对象

```go
func NewMysql(conf *Config) (*DB, error) {
	if conf.MaxIdleConns == 0 {
		conf.MaxIdleConns = 100
	}
	if conf.MaxOpenConns == 0 {
		conf.MaxOpenConns = 10
	}
	if conf.MaxLifetime == 0 {
		conf.MaxLifetime = 360
	}
	db, err := gorm.Open(mysql.Open(conf.DSN), &gorm.Config{
		SkipDefaultTransaction:                   true,
		NamingStrategy:                           nil,
		FullSaveAssociations:                     false,
		Logger:                                   &ormLog{},
		NowFunc:                                  nil,
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		IgnoreRelationshipsWhenMigrating:         false,
		DisableNestedTransaction:                 false,
		AllowGlobalUpdate:                        false,
		QueryFields:                              false,
		CreateBatchSize:                          0,
		TranslateError:                           false,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  nil,
	})

	if err != nil {
		return nil, err
	}
	sdb, err := db.DB()
	if err != nil {
		return nil, err
	}
	sdb.SetMaxIdleConns(conf.MaxIdleConns)
	sdb.SetMaxOpenConns(conf.MaxOpenConns)
	sdb.SetConnMaxLifetime(time.Second * time.Duration(conf.MaxLifetime))
    
    err = db.Use(NewCustomePlugin())
	if err != nil {
		return nil, err
	}
    
	return &DB{DB: db}, nil
}

func MustNewMysql(conf *Config) *DB {
	db, err := NewMysql(conf)
	if err != nil {
		panic(err)
	}
	return db
}

```

###### 自定义插件

```go
err = db.Use(NewCustomePlugin())
```

##### Metrics

##### Trace

##### 关注服务核心接口



##### Prometheus

下载地址: https://prometheus.io/download/

```yaml
# 配置解压后目录下的prometheus.yml
# vim prometheus.yml

......

scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: 'prometheus'

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs:
    - targets: ['localhost:9090']

  - job_name: 'file_ds'
    file_sd_configs:
    - files:
      - targets.json
```

```json
// touch target.json

// target.json配置

[
    {
        "targets": ["127.0.0.1:9101"],
        "labels": {
            "job": "follow",
            "app": "follow",
            "env": "test",
            "instance": "127.0.0.1:9101"
        }
    }
]
```



```shell
# 启动
./prometheus --config.file=prometheus.yml

# 浏览器访问
http://127.0.0.1:9090/
```

##### Jaeger

下载地址：https://www.jaegertracing.io/download/

```shell
# 启动
./jaeger-all-in-one

# 浏览器访问
http://127.0.0.1:16686/search
```

##### 测试

测试前需要启动以下服务

- Mysql
- Redis
- Etcd
- Prometheus
- Jaeger

###### 1.调用Follow进行关注

```json
// 参数
{
   "userId": 1,
   "followedUserId": 2  
}
```

- 查看结果
    ```sql
    -- 数据库
    select * from follow
    select * from follow_count;

    -- 缓存
    -- 无
    ```

- 查看日志
    ```json
    {"@timestamp":"2023-10-22T11:47:38.487+08:00","caller":"orm/orm.go:57","content":"[6.847ms] [rows:1] INSERT INTO `follow` (`user_id`,`followed_user_id`,`follow_status`,`create_time`,`update_time`) VALUES (1,2,1,'2023-10-22 11:47:38.48','2023-10-22 11:47:38.48')","duration":"6.8ms","level":"info","span":"1e40d3b77759bf02","trace":"773eb33b5ef1356ca176c0e45b97b695"}
    ```

    ```json
    {"@timestamp":"2023-10-22T11:47:38.493+08:00","caller":"orm/orm.go:57","content":"[5.938ms] [rows:1] INSERT INTO follow_count (user_id, follow_count) VALUES (1, 1) ON DUPLICATE KEY UPDATE follow_count = follow_count + 1","duration":"5.9ms","level":"info","span":"1e40d3b77759bf02","trace":"773eb33b5ef1356ca176c0e45b97b695"}

    {"@timestamp":"2023-10-22T11:47:38.493+08:00","caller":"orm/orm.go:57","content":"[0.438ms] [rows:1] INSERT INTO follow_count (user_id, fans_count) VALUES (2, 1) ON DUPLICATE KEY UPDATE fans_count = fans_count + 1","duration":"0.4ms","level":"info","span":"1e40d3b77759bf02","trace":"773eb33b5ef1356ca176c0e45b97b695"}
    ```

- 查看指标

    浏览器访问: http://127.0.0.1:9101/metrics

- 查看Trace

    浏览器访问: http://127.0.0.1:16686/search

###### 2. 调用UnFollow取消关注

```json
// 参数
{
    "userId": 1,
    "followedUserId": 2
}
```

- 查看结果

  ```sql
  -- 数据库
  select * from follow
  select * from follow_count;
  
  -- 缓存
  -- 无
  ```

- 查看日志

  ```json
  {"@timestamp":"2023-10-22T12:01:26.640+08:00","caller":"orm/orm.go:57","content":"[5.853ms] [rows:1] SELECT * FROM `follow` WHERE user_id = 1 AND followed_user_id = 2 ORDER BY `follow`.`id` LIMIT 1","duration":"5.9ms","level":"info","span":"5720f03000342ab4","trace":"34a819c091d483251f28030a6e5290ab"}
  {"@timestamp":"2023-10-22T12:01:26.644+08:00","caller":"orm/orm.go:57","content":"[3.496ms] [rows:1] UPDATE `follow` SET `follow_status`=2 WHERE id = 1","duration":"3.5ms","level":"info","span":"5720f03000342ab4","trace":"34a819c091d483251f28030a6e5290ab"}
  {"@timestamp":"2023-10-22T12:01:26.648+08:00","caller":"orm/orm.go:57","content":"[4.460ms] [rows:1] UPDATE follow_count SET follow_count = follow_count - 1 WHERE user_id = 1 AND follow_count \u003e 0","duration":"4.5ms","level":"info","span":"5720f03000342ab4","trace":"34a819c091d483251f28030a6e5290ab"}
  {"@timestamp":"2023-10-22T12:01:26.650+08:00","caller":"orm/orm.go:57","content":"[1.305ms] [rows:1] UPDATE follow_count SET fans_count = fans_count - 1 WHERE user_id = 2 AND fans_count \u003e 0","duration":"1.3ms","level":"info","span":"5720f03000342ab4","trace":"34a819c091d483251f28030a6e5290ab"}
  ```

- 查看指标

  

- 查看Trace

###### 3. 调用FollowList关注获取关注列表

- 增加关注信息

  调用Follow关注

  ```json
  {
      "userId": 1,
      "followedUserId": 2
  }
  ```

- 参数

  ```json
  {
      "userId": 1
  }
  ```

- 无缓存查看Trace

  ```json
  {"@timestamp":"2023-10-22T12:33:41.119+08:00","caller":"serverinterceptors/statinterceptor.go:90","content":"127.0.0.1:59747 - /service.Follow/FollowList - {\"userId\":1,\"cursor\":1697949221,\"pageSize\":20}","duration":"5.8ms","level":"info","span":"6ddc8012f901b120","trace":"6cd3304685ef474dfa3ab1c0e01b6d6d"}
  ```

- 有缓存查看Trace

  ```json
  {"@timestamp":"2023-10-22T12:35:43.519+08:00","caller":"serverinterceptors/statinterceptor.go:90","content":"127.0.0.1:59747 - /service.Follow/FollowList - {\"userId\":1,\"cursor\":1697949343,\"pageSize\":20}","duration":"4.3ms","level":"info","span":"da7518fa515ed257","trace":"37edfac1cc9de7d10b79cff2bf416430"}
  ```

  



----

### 第十一课

文档地址: https://pwmzlkcu3p.feishu.cn/docx/NYyHdpSzhoB8Zkxdb6NcpUGynH4
视频地址: https://www.bilibili.com/video/BV11u4y1Y7GC/

##### ElasticSearch

elasticsearch下载链接：https://www.elastic.co/guide/en/elasticsearch/reference/current/install-elasticsearch.html

kibana下载链接：https://www.elastic.co/guide/en/kibana/current/install.html

```shell
# 解压
tar -zxvf elasticsearch-8.10.4-darwin-x86_64.tar.gz -C /usr/local
tar -zxvf kibana-8.10.4-darwin-x86_64.tar.gz -C /usr/local

# 启动
cd /usr/local/elasticsearch-8.10.4
./bin/elasticsearch

/usr/local/kibana-8.10.4
./bin/kibana
```

验证是否成功

1. 浏览器输入：https://127.0.0.1:9200/

2. 弹出窗口：(注意看控制台上输出的内容)输入用户名和密码，用户名为elastic，默认密码在es启动后会输出在终端，见上图，也可以自定义

   ```sh
   Password for the elastic user (reset with `bin/elasticsearch-reset-password -u elastic`):
   iqmwv4TDc--S_PXF3yvT
   
   鈩癸笍  HTTP CA certificate SHA-256 fingerprint:
     df3ef3ade112405895c5541eb12506f725b92d74678083112ff1541f821b8abf
     
     鈩癸笍  Configure Kibana to use this cluster:
   鈥?Run Kibana and click the configuration link in the terminal when Kibana starts.
   鈥?Copy the following enrollment token and paste it into Kibana in your browser (valid for the next 30 minutes):
     eyJ2ZXIiOiI4LjE0LjAiLCJhZHIiOlsiMTkyLjE2OC4xNS4xMzM6OTIwMCJdLCJmZ3IiOiJkZjNlZjNhZGUxMTI0MDU4OTVjNTU0MWViMTI1MDZmNzI1YjkyZDc0Njc4MDgzMTEyZmYxNTQxZjgyMWI4YWJmIiwia2V5IjoibF9EbWdaQUJuckFUdG42QjRkRm46VEV6cl92UXRRSkdMaFRpTldqWEdDQSJ9
   
   ```

   

3. 浏览器输出如下即表示安装启动成功

   ```JSON
   {
       "name": "zhoushuangdeMBP",
       "cluster_name": "elasticsearch",
       "cluster_uuid": "L3_YeIFCRzGi7txrjVLdlQ",
       "version": {
           "number": "8.10.4",
           "build_flavor": "default",
           "build_type": "tar",
           "build_hash": "b4a62ac808e886ff032700c391f45f1408b2538c",
           "build_date": "2023-10-11T22:04:35.506990650Z",
           "build_snapshot": false,
           "lucene_version": "9.7.0",
           "minimum_wire_compatibility_version": "7.17.0",
           "minimum_index_compatibility_version": "7.0.0"
       },
       "tagline": "You Know, for Search"
   }
   ```

4. 复制输出 http://localhost:5601/?code=815471地址到浏览器

5. 将上面es启动后在终端输出的秘钥拷贝进去然后点击`Configure Elastic`
### 第十二课

文档地址: https://pwmzlkcu3p.feishu.cn/docx/QZcKdB4VXoUCRDxGilfcTaw7n8e
视频地址: https://www.bilibili.com/video/BV1u64y177rL

### 第十三课

文档地址: https://pwmzlkcu3p.feishu.cn/docx/BDNUdhmP4oec6ix1P1ZcqF3rnIc
视频地址: https://www.bilibili.com/video/BV1bH4y1C7Uj

### 第十四课

文档地址: https://pwmzlkcu3p.feishu.cn/docx/Ydd4dG8OSobJ1rxJFgacOBKvnSg
视频地址: https://www.bilibili.com/video/BV1qK411b7cS