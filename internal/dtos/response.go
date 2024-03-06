package dtos

import (
	"med-chat-bot/internal/meta"
)

// Response response custom message
type Response struct {
	Meta meta.Meta   `json:"meta"`
	Data interface{} `json:"data" swaggertype:"object"`
}

type PaginationResponse struct {
	Meta           meta.Meta       `json:"meta"`
	PaginationInfo *PaginationInfo `json:"pagination"`
	Data           interface{}     `json:"data" swaggertype:"object"`
}

type PaginationInfo struct {
	PageSize   int64 `json:"pageSize"`
	PageOffset int64 `json:"pageOffset"`
}

type ListParam struct {
	PageOffset int64
	PageSize   int64
	Pagination bool
	Preload    bool
	OrderBy    *string
}

type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

const (
	QueryValueAll  = "*"
	QueryValueNone = "-"
)

const (
	DefaultPageSize int64 = 50
	MinPageSize     int64 = 10
	MaxPageSize     int64 = 1000
)
