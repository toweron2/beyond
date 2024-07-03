# 文章服务

```shell
goctl api go --dir=./ --api article.api --style go_zero
```


```shell
cd beyond/application/article/rpc

goctl rpc protoc ./article.proto --go_out=. --go-grpc_out=. --zrpc_out=./ --style go_zero
```

```shell
cd beyond/application/article

goctl model mysql datasource --dir ./model -t "article" --url "root:200212..@tcp(127.0.0.1:3306)/beyond_article" --style go_zero
goctl model mysql datasource --dir ./model -t "*" --url "root:200212..@tcp(127.0.0.1:3306)/beyond_article" -c true --style go_zero 
```

```shell
protoc -I=. --go_out=. --proto_path=/usr/local/include/ ./article.proto
protoc -I=. --go_out=. --go-grpc_out=. --proto_path=/usr/local/include/ ./article.proto

```

缓存击穿经常发生在热点数据过期失效的时候，那么我们不让缓存失效不就好了，每次查询缓存的时候使用Exists来判断key是否存在，如果存在就使用Expire给缓存续期，既然是热点数据通过不断地续期也就不会过期了

缓存穿透是指要访问的数据既不在缓存中，也不在数据库中，导致请求在访问缓存时，发生缓存缺失，再去访问数据库时，发现数据库中也没有要访问的数据。此时也就没办法从数据库中读出数据再写入缓存来服务后续的请求，类似的请求如果多的话就会给缓存和数据库带来巨大的压力。
针对缓存穿透问题，解决办法其实很简单，就是缓存一个空值，避免每次都透传到数据库，缓存的时间可以设置短一点，比如1分钟，其实上文已经有提到了，当我们访问不存在的数据的时候，go-zero框架会帮我们自动加上空缓存，比如我们访问id为999的文章，该文章在数据库中是不存在的。

缓存雪崩时指大量的的应用请求无法在Redis缓存中进行处理，紧接着应用将大量的请求发送到数据库，导致数据库被打挂，好惨呐！！缓存雪崩一般是由两个原因导致的，应对方案也不太一样。
第一个原因是：缓存中有大量的数据同时过期，导致大量的请求无法得到正常处理。
针对大量数据同时失效带来的缓存雪崩问题，一般的解决方案是要避免大量的数据设置相同的过期时间，如果业务上的确有要求数据要同时失效，那么可以在过期时间上加一个较小的随机数，这样不同的数据过期时间不同，但差别也不大，避免大量数据同时过期，也基本能满足业务的需求。
第二个原因是：Redis出现了宕机，没办法正常响应请求了，这就会导致大量请求直接打到数据库，从而发生雪崩
针对这类原因一般我们需要让我们的数据库支持熔断，让数据库压力比较大的时候就触发熔断，丢弃掉部分请求，当然熔断是对业务有损的。