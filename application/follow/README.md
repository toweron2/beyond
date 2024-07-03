



```shell
cd beyond/application/follow

goctl model mysql datasource --dir ./model -t "article" --url "root:200212..@tcp(127.0.0.1:3306)/beyond_article" --style go_zero
goctl model mysql datasource --dir ./model -t "*" --url "root:200212..@tcp(127.0.0.1:3306)/beyond_article" -c true --style go_zero 
```

```shell
cd beyond/application/article/rpc

goctl rpc protoc ./follow.proto --go_out=. --go-grpc_out=. --zrpc_out=./ --style go_zero
```

```shell
protoc -I=. --go_out=. --proto_path=/usr/local/include/ ./follow.proto
protoc -I=. --go_out=. --go-grpc_out=. --proto_path=/usr/local/include/ ./follow.proto

```