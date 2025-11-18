/*
Copyright © 2025 lixw
*/
package web

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ethanli-dev/go-app-layout/docs"
	"github.com/ethanli-dev/go-app-layout/pkg/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swagfiles "github.com/swaggo/files"
	ginswag "github.com/swaggo/gin-swagger"
)

type Options struct {
	address        string
	readTimeout    time.Duration
	writeTimeout   time.Duration
	idleTimeout    time.Duration
	maxHeaderBytes int
	basePath       string
	staticPath     string
	fs             fs.FS
	middleware     []gin.HandlerFunc
}

type Option func(*Options)

func WithAddress(address string) Option {
	return func(o *Options) {
		o.address = address
	}
}

func WithReadTimeout(readTimeout time.Duration) Option {
	return func(o *Options) {
		o.readTimeout = readTimeout
	}
}

func WithWriteTimeout(writeTimeout time.Duration) Option {
	return func(o *Options) {
		o.writeTimeout = writeTimeout
	}
}

func WithIdleTimeout(idleTimeout time.Duration) Option {
	return func(o *Options) {
		o.idleTimeout = idleTimeout
	}
}

func WithMaxHeaderBytes(maxHeaderBytes int) Option {
	return func(o *Options) {
		o.maxHeaderBytes = maxHeaderBytes
	}
}

func WithBasePath(basePath string) Option {
	return func(o *Options) {
		o.basePath = basePath
	}
}

// WithStaticFS 使用本地文件系统目录
func WithStaticFS(path, root string) Option {
	return func(o *Options) {
		o.fs = os.DirFS(root)
		o.staticPath = path
	}
}

// WithEmbedFS 使用嵌入式资源
func WithEmbedFS(path string, efs fs.FS) Option {
	return func(o *Options) {
		o.fs = efs
		o.staticPath = path
	}
}

func WithMiddleware(middleware ...gin.HandlerFunc) Option {
	return func(o *Options) {
		o.middleware = append(o.middleware, middleware...)
	}
}

type Server struct {
	httpSrv   *http.Server
	baseRoute *gin.RouterGroup
	engine    *gin.Engine
}

func New(options ...Option) *Server {
	opts := &Options{
		address:        ":8080",
		readTimeout:    5 * time.Second,
		writeTimeout:   10 * time.Second,
		idleTimeout:    30 * time.Second,
		maxHeaderBytes: 1 << 20,
	}

	for _, option := range options {
		option(opts)
	}

	docs.SwaggerInfo.BasePath = opts.basePath

	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.RedirectTrailingSlash = true
	engine.RedirectFixedPath = true
	// 中间件注册顺序（关键！）
	// 1. 跨域处理 → 2. 追踪ID → 3. 日志 → 4. 异常恢复
	engine.Use(
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"*"},
			ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
		middleware.RequestId(),
		middleware.Logger(),
		middleware.Recovery(),
		middleware.I18n(),
	)
	engine.Use(opts.middleware...)

	engine.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().UnixMilli()})
	})

	engine.GET("/swagger/*any", ginswag.WrapHandler(swagfiles.Handler))

	if opts.fs != nil {
		engine.StaticFS(opts.staticPath, http.FS(opts.fs))
		slog.Info("serving static files", "path", opts.staticPath)
	}
	
	return &Server{
		baseRoute: engine.Group(opts.basePath),
		engine:    engine,
		httpSrv: &http.Server{
			Addr:           opts.address,
			Handler:        engine,
			IdleTimeout:    opts.idleTimeout,
			ReadTimeout:    opts.readTimeout,
			WriteTimeout:   opts.writeTimeout,
			MaxHeaderBytes: opts.maxHeaderBytes,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	slog.InfoContext(ctx, "starting http server", "addr", s.httpSrv.Addr)
	startErrCh := make(chan error, 1)
	go func() {
		if err := s.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			startErrCh <- err
		} else {
			close(startErrCh) // 正常关闭，不产生错误
		}
	}()
	select {
	case err := <-startErrCh:
		return fmt.Errorf("http server failed to start: %w", err)
	case <-time.After(100 * time.Millisecond):
		// 短暂延迟，基本能判断端口绑定是否成功
		slog.InfoContext(ctx, "http server started successfully")
		return nil
	}
}

func (s *Server) Stop(ctx context.Context) error {
	slog.InfoContext(ctx, "stopping http server", "addr", s.httpSrv.Addr)
	if err := s.httpSrv.Shutdown(ctx); err != nil {
		// 若优雅关闭失败，尝试强制关闭
		if closeErr := s.httpSrv.Close(); closeErr != nil {
			slog.ErrorContext(ctx, "failed to force close http server", "err", closeErr)
		}
		return fmt.Errorf("shutdown http server failed: %w", err)
	}
	slog.InfoContext(ctx, "http server shutdown complete")
	return nil
}

func (s *Server) UseRoutes(routeFuncs ...func(*gin.RouterGroup)) *Server {
	for _, fn := range routeFuncs {
		fn(s.baseRoute)
	}
	return s
}
