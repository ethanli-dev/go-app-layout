/*
Copyright © 2025 lixw
*/
package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 判断是否为连接中断错误
				brokenPipe := isBrokenPipeError(err)

				// 生成请求信息快照
				reqDump, dumpErr := httputil.DumpRequest(ctx.Request, false)
				reqStr := string(reqDump)
				if dumpErr != nil {
					reqStr = fmt.Sprintf("[failed to dump request: %v]", dumpErr)
				}

				// 获取堆栈跟踪信息
				stack := getPanicStack(3)

				// 准备结构化日志字段
				attrs := []any{
					slog.String("method", ctx.Request.Method),
					slog.String("url", ctx.Request.URL.String()),
					slog.String("request", reqStr),
					slog.String("stack", stack),
					slog.String("error_type", map[bool]string{true: "broken_pipe", false: "panic"}[brokenPipe]),
				}

				// 处理非error类型的panic
				errMsg := fmt.Sprintf("%v", err)
				if errObj, ok := err.(error); ok {
					errMsg = errObj.Error()
				}
				attrs = append(attrs, slog.String("error", errMsg))

				// 记录结构化错误日志
				slog.ErrorContext(ctx, "recovered from panic", attrs...)

				// 处理连接中断错误
				if brokenPipe {
					// 对于broken pipe，直接关闭连接，避免进一步写入
					_ = ctx.Error(fmt.Errorf("broken pipe: %v", err))
					ctx.Abort()
					return
				}

				// 其他错误返回500
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		ctx.Next()
	}
}

// isBrokenPipeError 判断错误是否为"broken pipe"或"connection reset by peer"（客户端断开连接）
func isBrokenPipeError(err interface{}) bool {
	// 安全转换panic值为error（处理非error类型的panic）
	errObj, ok := err.(error)
	if !ok {
		return false
	}

	// 检查是否为网络操作错误
	var netOpErr *net.OpError
	if !errors.As(errObj, &netOpErr) {
		return false
	}

	// 检查是否为系统调用错误
	var syscallErr *os.SyscallError
	if errors.As(netOpErr.Err, &syscallErr) {
		errStr := strings.ToLower(syscallErr.Error())
		return strings.Contains(errStr, "broken pipe") ||
			strings.Contains(errStr, "connection reset by peer")
	}

	// 检查是否为关闭的连接错误
	if errors.Is(netOpErr.Err, os.ErrClosed) {
		return true
	}

	return false
}

// getPanicStack 获取panic发生时的堆栈跟踪信息（跳过指定数量的调用帧）
func getPanicStack(skip int) string {
	var pcs [32]uintptr
	// 增加1个跳过帧，避免包含当前函数自身
	n := runtime.Callers(skip+1, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var stack strings.Builder
	stack.WriteString("stack trace:\n") // 添加堆栈标记
	for i := 0; ; i++ {
		frame, more := frames.Next()
		// 格式化每一行堆栈信息，添加序号便于阅读
		stack.WriteString(fmt.Sprintf("#%d %s\n    %s:%d\n",
			i, frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}
	return stack.String()
}
