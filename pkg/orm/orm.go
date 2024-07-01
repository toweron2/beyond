package orm

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type (
	Config struct {
		DSN          string
		MaxOpenConns int
		MaxIdleConns int
		MaxLifetime  int
	}
	DB struct {
		*gorm.DB
	}
	ormLog struct {
		LogLevel logger.LogLevel
	}
)

func (l *ormLog) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

func (l *ormLog) Info(ctx context.Context, format string, v ...interface{}) {
	if l.LogLevel < logger.Info {
		return
	}
	logx.ErrorLevel
	logx.WithContext(ctx).Infof(format, v...)
}

func (l *ormLog) Warn(ctx context.Context, s string, i ...interface{}) {
	// TODO implement me
	panic("implement me")
}

func (l *ormLog) Error(ctx context.Context, s string, i ...interface{}) {
	// TODO implement me
	panic("implement me")
}

func (l *ormLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// TODO implement me
	panic("implement me")
}
