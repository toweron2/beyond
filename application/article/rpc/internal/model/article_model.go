package model

import (
	"context"
	"database/sql"
	"fmt"
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
		ArticlesByUserId(ctx context.Context, uid, likeNum, limit int64, status int32, pubTime, sortField string) ([]*Article, error)
		UpdateArticleStatus(ctx context.Context, id int64, status int32) error
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

func (c *customArticleModel) ArticlesByUserId(ctx context.Context, uid, likeNum, limit int64, status int32, pubTime, sortField string) ([]*Article, error) {
	where := map[string]any{
		"_orderby": sortField,
		"_limit":   limit,
	}
	where["user_id"] = uid
	where["status"] = status

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

func (c *customArticleModel) UpdateArticleStatus(ctx context.Context, id int64, status int32) error {
	beyondArticleArticleIdKey := fmt.Sprintf("%s%v", cacheBeyondArticleArticleIdPrefix, id)
	_, err := c.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		query := fmt.Sprintf("update %s set status = ? where `id` = ?", c.table)
		return conn.ExecCtx(ctx, query, status, id)
	}, beyondArticleArticleIdKey)
	return err
}
