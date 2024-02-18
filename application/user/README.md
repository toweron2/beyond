# 用户服务

goctl rpc protoc ./user.proto --go_out=. --go-grpc_out=. --zrpc_out=./ --style go_zero

goctl model mysql datasource --dir ./internal/model --table user --cache true --url "root:200212..@tcp(127.0.0.1:3306)/beyond_user"