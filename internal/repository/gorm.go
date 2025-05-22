package repository

import (
	"context"

	"gorm.io/gorm"
)

// GormDB GORM 数据库适配器
type GormDB struct {
	*gorm.DB
}

// NewGormDB 创建 GORM 数据库适配器
func NewGormDB(db *gorm.DB) DB {
	return &GormDB{DB: db}
}

// WithContext 实现 DB 接口
func (db *GormDB) WithContext(ctx context.Context) DB {
	return &GormDB{DB: db.DB.WithContext(ctx)}
}

// Create 实现 DB 接口
func (db *GormDB) Create(value interface{}) error {
	return db.DB.Create(value).Error
}

// Save 实现 DB 接口
func (db *GormDB) Save(value interface{}) error {
	return db.DB.Save(value).Error
}

// Delete 实现 DB 接口
func (db *GormDB) Delete(value interface{}, conds ...interface{}) error {
	return db.DB.Delete(value, conds...).Error
}

// First 实现 DB 接口
func (db *GormDB) First(dest interface{}, conds ...interface{}) error {
	return db.DB.First(dest, conds...).Error
}

// Find 实现 DB 接口
func (db *GormDB) Find(dest interface{}, conds ...interface{}) error {
	return db.DB.Find(dest, conds...).Error
}

// Where 实现 DB 接口
func (db *GormDB) Where(query interface{}, args ...interface{}) DB {
	return &GormDB{DB: db.DB.Where(query, args...)}
}

// Model 实现 DB 接口
func (db *GormDB) Model(value interface{}) DB {
	return &GormDB{DB: db.DB.Model(value)}
}

// Table 实现 DB 接口
func (db *GormDB) Table(name string) DB {
	return &GormDB{DB: db.DB.Table(name)}
}

// Exec 实现 DB 接口
func (db *GormDB) Exec(sql string, values ...interface{}) error {
	return db.DB.Exec(sql, values...).Error
}

// Transaction 实现 DB 接口
func (db *GormDB) Transaction(fc func(tx DB) error) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		return fc(&GormDB{DB: tx})
	})
}

// AutoMigrate 实现 DB 接口
func (db *GormDB) AutoMigrate(dst ...interface{}) error {
	return db.DB.AutoMigrate(dst...)
}

// Count 实现 DB 接口
func (db *GormDB) Count(value *int64) error {
	result := db.DB.Count(value)
	return result.Error
}

// Offset 实现 DB 接口
func (db *GormDB) Offset(offset int) DB {
	return &GormDB{DB: db.DB.Offset(offset)}
}

// Limit 实现 DB 接口
func (db *GormDB) Limit(limit int) DB {
	return &GormDB{DB: db.DB.Limit(limit)}
}
