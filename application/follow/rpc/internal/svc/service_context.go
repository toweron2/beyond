package svc

import (
	"beyond/application/follow/rpc/internal/config"
	"beyond/application/follow/rpc/internal/model"
	"beyond/pkg/orm"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config           config.Config
	DB               *orm.DB
	FollowModel      *model.FollowModel
	FollowCountModel *model.FollowCountModel
	BizRedis         *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := orm.MustNewMysql(&c.DB)

	rds := redis.MustNewRedis(c.BizRedis)
	return &ServiceContext{
		Config:   c,
		DB:       db,
		BizRedis: rds,
	}
}
