# 文章服务

```shell
goctl api go --dir=./ --api article.api --style go_zero
```


```shell
cd beyond/application/article/rpc

goctl rpc protoc ./article.proto --go_out=. --go-grpc_out=. --zrpc_out=./ --style go_zero
```

```shell
cd beyond/application/article/rpc

goctl model mysql datasource --dir ./internal/model -t "article" --url "root:200212..@tcp(127.0.0.1:3306)/beyond_article" --style go_zero
goctl model mysql datasource --dir ./internal/model -t "*" --url "root:200212..@tcp(127.0.0.1:3306)/beyond_article" -c true --style go_zero 
```

```shell
protoc -I=. --go_out=. --proto_path=/usr/local/include/ ./article.proto
protoc -I=. --go_out=. --go-grpc_out=. --proto_path=/usr/local/include/ ./article.proto

```