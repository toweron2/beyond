package svc

import (
	"beyond/application/article/internal/model"
	"beyond/application/article/mq/internal/config"
	logic2 "beyond/application/article/mq/internal/logic"

	"beyond/application/user/rpc/user"
	"beyond/pkg/es"
	"context"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	ArticleModel model.ArticleModel
	BizRedis     *redis.Redis

	UserRpc user.User
	Es      *es.Es
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdf := redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	}
	conn := sqlx.NewMysql(c.DataSource)

	return &ServiceContext{
		Config:       c,
		ArticleModel: model.NewArticleModel(conn, cache.CacheConf{{rdf, 100}}),
		BizRedis:     redis.MustNewRedis(rdf),
		UserRpc:      user.NewUser(zrpc.MustNewClient(c.UserRPC)),
		Es:           es.MustNewEs(&c.Es),
	}
}

func Consumers(ctx context.Context, svcCtx *ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConsumerConf, logic2.NewArticleLikeNumLogic(ctx, svcCtx)),
		kq.MustNewQueue(svcCtx.Config.ArticleKqConsumerConf, logic2.NewArticleLogic(ctx, svcCtx)),
	}
}
