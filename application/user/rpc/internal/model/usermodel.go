package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		FindByMobile(ctx context.Context, mobile string) (*User, error)
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c, opts...),
	}
}

func (m *customUserModel) FindByMobile(ctx context.Context, mobile string) (*User, error) {
	var user User
	err := m.QueryRowNoCacheCtx(ctx, &user, fmt.Sprintf("select %s from %s where `mobile` = ? limit 1", userRows, m.table), mobile)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	switch err {
	case nil:
		return &user, nil
	case sqlc.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}
