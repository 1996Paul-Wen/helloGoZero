package util

import (
	"context"

	"github.com/google/uuid"
)

// contextKey 是自定义类型，避免 key 冲突
type contextKey string

const TRACEKEY = contextKey("trace-id")

func LoadTraceFrom(ctx context.Context) string {
	untypeTrace := ctx.Value(TRACEKEY)
	if untypeTrace != nil {
		return untypeTrace.(string)
	}
	return ""
}

func NewTraceID() string {
	return uuid.New().String()
}
