package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	ArticleRPC zrpc.RpcClientConf
	UserRPC    zrpc.RpcClientConf

	/*Oss struct {
		Endpoint         string
		AccessKeyId      string
		AccessKeySecret  string
		BucketName       string
		ConnectTimeout   int64 `json:",optional"`
		ReadWriteTimeout int64 `json:",optional"`
	}*/
	Qiniu struct {
		AccessKey string
		SecretKey string `json:",optional"`
		Bucket    string // 命名空间
		// 机房位置
		// Zone          string `json:",optional"`
		// ImgPath       string
		UseHTTPS      bool `json:",default=false"`
		UseCdnDomains bool `json:",default=true"`
	}
}
