package repository

import (
	"context"
)

// DB 数据库接口
type DB interface {
	WithContext(ctx context.Context) DB
	Create(value interface{}) error
	Save(value interface{}) error
	Delete(value interface{}, conds ...interface{}) error
	First(dest interface{}, conds ...interface{}) error
	Find(dest interface{}, conds ...interface{}) error
	Where(query interface{}, args ...interface{}) DB
	Model(value interface{}) DB
	Table(name string) DB
	Exec(sql string, values ...interface{}) error
	Transaction(fc func(tx DB) error) error
	AutoMigrate(dst ...interface{}) error
	Count(value *int64) error
	Offset(offset int) DB
	Limit(limit int) DB
}
