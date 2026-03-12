package handler

import (
	"bytes"
	"context"
	"net/http"

	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/util"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

// 日志中间件
func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logx.WithContext(r.Context())
		logger.Infof("Request: %s %s. trace: %s", r.Method, r.RequestURI, util.LoadTraceFrom(r.Context()))

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // 默认 200
			body:           &bytes.Buffer{},
		}
		next(rw, r) // 调用下一个中间件或最终处理函数
		// todo 打印response
		logger.Infof("Response: status=%d, body=%s", rw.statusCode, rw.body.String())

	}
}

// responseWriter 自定义 ResponseWriter，用于捕获状态码和响应体
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

// WriteHeader 捕获状态码
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write 捕获响应数据
func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b) // 缓存 body，用于日志记录
	return rw.ResponseWriter.Write(b)
}

// 鉴权中间件
// func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		token := r.Header.Get("Authorization")
// 		if token != "valid-token" {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}
// 		next(w, r)
// 	}
// }

// TraceMiddleware 生成或透传 Trace ID 并注入 context
func TraceMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. 尝试从请求头获取 trace ID
		traceID := r.Header.Get("X-Trace-ID")
		if traceID == "" {
			// 2. 如果不存在，生成新的
			traceID = uuid.NewString()
		}

		// 3. 将 trace ID 存入 context
		ctx := context.WithValue(r.Context(), util.TRACEKEY, traceID)

		// 4. 调用下一个 handler，使用带 trace ID 的 context
		next(w, r.WithContext(ctx))
	}
}
