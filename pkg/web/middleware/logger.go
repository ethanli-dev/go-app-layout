/*
Copyright © 2025 lixw
*/
package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/ethanli-dev/go-app-layout/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	requestBodyLimit   = 1 << 20
	HeaderKeyRequestID = "X-Request-Id"
)

var (
	bufPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 4096))
		},
	}
)

func RequestId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodOptions {
			ctx.Next()
			return
		}
		requestID := ctx.GetHeader(HeaderKeyRequestID)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx.Set(logging.ContextKeyTraceID, requestID)
		ctx.Writer.Header().Set(HeaderKeyRequestID, requestID)
		ctx.Next()
	}
}

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodOptions {
			ctx.Next()
			return
		}
		start := time.Now()
		var (
			reqBody         []byte
			isBodyTruncated bool
		)
		if ctx.Request.Body != http.NoBody {
			buf := bufPool.Get().(*bytes.Buffer)
			defer func() {
				buf.Reset()
				bufPool.Put(buf)
			}()
			// 限制读取的最大字节数1MB，避免内存溢出
			limitedReader := &io.LimitedReader{
				R: ctx.Request.Body,
				N: requestBodyLimit,
			}
			if n, err := io.Copy(buf, limitedReader); err == nil {
				// 检查是否被截断
				isBodyTruncated = limitedReader.N <= 0 && n >= requestBodyLimit
				reqBody = buf.Bytes()
				// 重置请求体供后续处理使用
				ctx.Request.Body = io.NopCloser(bytes.NewReader(reqBody))
			}
		}

		defer func() {
			cost := time.Since(start)
			status := ctx.Writer.Status()
			bodyStr := string(reqBody)
			if len(reqBody) > 0 && isBodyTruncated {
				bodyStr += " [truncated]"
			}

			attrs := []any{
				slog.String("method", ctx.Request.Method),
				slog.String("url", ctx.Request.URL.String()),
				slog.Int("status", status),
				slog.Duration("cost", cost),
				slog.String("request", bodyStr),
			}
			switch {
			case status >= http.StatusInternalServerError:
				slog.ErrorContext(ctx, "request completed with error", attrs...)
			case status >= http.StatusBadRequest:
				slog.WarnContext(ctx, "request completed with warning", attrs...)
			default:
				slog.InfoContext(ctx, "request completed successfully", attrs...)
			}
		}()
		ctx.Next()
	}
}
