package pagination

import (
	"ai-course/internal/repository"
	"context"
)

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// Params 分页参数
type Params struct {
	Page     int `form:"page" binding:"required,min=1"`     // 页码
	PageSize int `form:"page_size" binding:"min=1,max=100"` // 每页数量
}

// Response 分页响应
type Response struct {
	Total    int64       `json:"total"`     // 总记录数
	Page     int         `json:"page"`      // 当前页码
	PageSize int         `json:"page_size"` // 每页数量
	List     interface{} `json:"list"`      // 数据列表
}

// Paginate 执行分页查询
func Paginate[T any](ctx context.Context, db repository.DB, params *Params, model interface{}, result *[]T) (*Response, error) {
	// 设置默认分页大小
	if params.PageSize == 0 {
		params.PageSize = DefaultPageSize
	}
	if params.PageSize > MaxPageSize {
		params.PageSize = MaxPageSize
	}

	var total int64
	offset := (params.Page - 1) * params.PageSize

	// 获取总记录数
	if err := db.WithContext(ctx).Model(model).Count(&total); err != nil {
		return nil, err
	}

	// 获取分页数据
	if err := db.WithContext(ctx).Offset(offset).Limit(params.PageSize).Find(result); err != nil {
		return nil, err
	}

	return &Response{
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
		List:     result,
	}, nil
}

// ValidateAndSetDefaults 验证并设置默认值
func ValidateAndSetDefaults(params *Params) {
	if params.PageSize == 0 {
		params.PageSize = DefaultPageSize
	}
	if params.PageSize > MaxPageSize {
		params.PageSize = MaxPageSize
	}
	if params.Page < 1 {
		params.Page = 1
	}
}
