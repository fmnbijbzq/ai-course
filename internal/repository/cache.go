package repository

import (
	"context"
	"time"
)

// Cache 缓存接口
type Cache interface {
	// Get 获取缓存
	Get(ctx context.Context, key string) (interface{}, error)
	// Set 设置缓存
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// Delete 删除缓存
	Delete(ctx context.Context, key string) error
}

// NoOpCache 空缓存实现，用于不需要缓存时
type NoOpCache struct{}

// NewNoOpCache 创建空缓存实例
func NewNoOpCache() Cache {
	return &NoOpCache{}
}

// Get 获取缓存
func (c *NoOpCache) Get(ctx context.Context, key string) (interface{}, error) {
	return nil, nil
}

// Set 设置缓存
func (c *NoOpCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return nil
}

// Delete 删除缓存
func (c *NoOpCache) Delete(ctx context.Context, key string) error {
	return nil
}
