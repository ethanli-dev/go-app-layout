/*
Copyright © 2025 lixw
*/
package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Service interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type ServiceFunc struct {
	StartFunc func(context.Context) error
	StopFunc  func(context.Context) error
}

func (s ServiceFunc) Start(ctx context.Context) error {
	if s.StartFunc == nil {
		return nil
	}
	return s.StartFunc(ctx)
}

func (s ServiceFunc) Stop(ctx context.Context) error {
	if s.StopFunc == nil {
		return nil
	}
	return s.StopFunc(ctx)
}

type App struct {
	services        []Service
	startTimeout    time.Duration
	shutdownTimeout time.Duration
}

type Option func(*App)

func WithStartTimeout(startTimeout time.Duration) Option {
	return func(a *App) {
		a.startTimeout = startTimeout
	}
}

func WithShutdownTimeout(shutdownTimeout time.Duration) Option {
	return func(a *App) {
		a.shutdownTimeout = shutdownTimeout
	}
}

func New(options ...Option) *App {
	app := &App{
		startTimeout:    15 * time.Second,
		shutdownTimeout: 15 * time.Second,
	}
	for _, option := range options {
		option(app)
	}
	return app
}

func (a *App) Use(services ...Service) *App {
	a.services = append(a.services, services...)
	return a
}

func (a *App) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "starting application")
	startCtx, cancel := context.WithTimeout(ctx, a.startTimeout)
	defer cancel()
	if err := withTimeout(startCtx, a.doStart); err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}
	slog.InfoContext(ctx, "application started successfully")

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	select {
	case <-quit:
		slog.InfoContext(ctx, "received interrupt signal, shutting down")
	case <-ctx.Done():
		slog.InfoContext(ctx, "context cancelled, shutting down")
	}
	slog.InfoContext(ctx, "shutting down application")

	// 优雅关闭
	stopCtx, cancel := context.WithTimeout(ctx, a.shutdownTimeout)
	defer cancel()

	if err := withTimeout(stopCtx, a.doStop); err != nil {
		return fmt.Errorf("failed to shutdown application: %w", err)
	}

	slog.InfoContext(ctx, "application shutdown complete")
	return nil
}

func (a *App) doStart(ctx context.Context) error {
	var started []Service
	for _, event := range a.services {
		// 检查上下文是否已取消
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := event.Start(ctx); err != nil {
			if e := a.doRelease(ctx, started); e != nil {
				return errors.Join(err, e)
			}
			return err
		}
		started = append(started, event)
	}
	return nil
}

func (a *App) doStop(ctx context.Context) error {
	return a.doRelease(ctx, a.services)
}

func (a *App) doRelease(ctx context.Context, services []Service) error {
	var errs []error

	// 逆序停止已启动的事件，确保资源释放顺序正确
	for i := len(services) - 1; i >= 0; i-- {
		// 检查上下文是否已取消
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := services[i].Stop(ctx); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func withTimeout(ctx context.Context, callback func(context.Context) error) error {
	errChan := make(chan error, 1)
	// 创建一个用于通知goroutine退出的通道
	done := make(chan struct{})
	defer close(done) // 确保在函数退出时关闭通道，通知goroutine可以退出
	go func() {
		callbackExited := false
		defer func() {
			if !callbackExited {
				errChan <- errors.New("goroutine exited without returning")
			}
		}()

		// 执行回调并捕获错误
		errChan <- callback(ctx)
		callbackExited = true
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		// 上下文已结束时优先返回超时错误
		return err
	}
}
