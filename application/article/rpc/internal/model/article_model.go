package model

import (
	"context"
	sqlBuilder "github.com/didi/gendry/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ArticleModel = (*customArticleModel)(nil)

type (
	// ArticleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customArticleModel.
	ArticleModel interface {
		articleModel
		ArticleByUserId(ctx context.Context, uid, likeNum, limit int64, pubTime, sortField string) ([]*Article, error)
	}

	customArticleModel struct {
		*defaultArticleModel
	}
)

// NewArticleModel returns a model for the database table.
func NewArticleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ArticleModel {
	return &customArticleModel{
		defaultArticleModel: newArticleModel(conn, c, opts...),
	}
}

func (c *customArticleModel) ArticleByUserId(ctx context.Context, uid, likeNum, limit int64, pubTime, sortField string) ([]*Article, error) {
	where := map[string]any{
		"_orderby": sortField,
		"_limit":   limit,
	}
	where["user_id"] = uid

	if sortField == "like_num" {
		where["like_num <"] = likeNum
	} else {
		where["publish_time <"] = pubTime
	}

	cond, vals, err := sqlBuilder.BuildSelect(c.table, where, articleFieldNames)
	if err != nil {
		return nil, err
	}
	var articles []*Article
	err = c.QueryRowsNoCacheCtx(ctx, articles, cond, vals...)
	if err != nil {
		return nil, err
	}

	return articles, err
}
