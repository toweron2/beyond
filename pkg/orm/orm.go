package orm

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _ logger.Interface = (*ormLog)(nil)

type (
	Config struct {
		DSN          string
		MaxOpenConns int `json:",default=10"`
		MaxIdleConns int `json:",default=100"`
		MaxLifetime  int `json:",default=3600"`
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
	if l.LogLevel <= logger.Info {
		return
	}
	logx.WithContext(ctx).Infof(format, v...)
}

func (l *ormLog) Warn(ctx context.Context, format string, v ...interface{}) {
	logx.WithContext(ctx).Infof(format, v...)
}

func (l *ormLog) Error(ctx context.Context, format string, v ...interface{}) {
	logx.WithContext(ctx).Errorf(format, v...)
}

func (l *ormLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	logx.WithContext(ctx).WithDuration(elapsed).Infof("[%.3fms] [rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
}

func NewMysql(conf *Config) (*DB, error) {
	if conf.MaxIdleConns == 0 {
		conf.MaxIdleConns = 100
	}
	if conf.MaxOpenConns == 0 {
		conf.MaxOpenConns = 10
	}
	if conf.MaxLifetime == 0 {
		conf.MaxLifetime = 360
	}
	db, err := gorm.Open(mysql.Open(conf.DSN), &gorm.Config{
		SkipDefaultTransaction:                   true,
		NamingStrategy:                           nil,
		FullSaveAssociations:                     false,
		Logger:                                   &ormLog{},
		NowFunc:                                  nil,
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		IgnoreRelationshipsWhenMigrating:         false,
		DisableNestedTransaction:                 false,
		AllowGlobalUpdate:                        false,
		QueryFields:                              false,
		CreateBatchSize:                          0,
		TranslateError:                           false,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  nil,
	})

	if err != nil {
		return nil, err
	}
	sdb, err := db.DB()
	if err != nil {
		return nil, err
	}
	sdb.SetMaxIdleConns(conf.MaxIdleConns)
	sdb.SetMaxOpenConns(conf.MaxOpenConns)
	sdb.SetConnMaxLifetime(time.Second * time.Duration(conf.MaxLifetime))

	err = db.Use(NewCustomePlugin())
	if err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil
}

func MustNewMysql(conf *Config) *DB {
	db, err := NewMysql(conf)
	if err != nil {
		panic(err)
	}
	return db
}
