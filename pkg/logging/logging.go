/*
Copyright © 2025 lixw
*/
package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

const ContextKeyTraceID = "traceId"

type TraceContextHandler struct {
	slog.Handler
}

func (h *TraceContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx != nil {
		if traceId, ok := ctx.Value(ContextKeyTraceID).(string); ok {
			r.Add(ContextKeyTraceID, slog.StringValue(traceId))
		}
	}
	return h.Handler.Handle(ctx, r)
}

type Options struct {
	level        string
	path         string
	maxAge       int
	maxSize      int
	maxBackups   int
	compress     bool
	format       string
	enableStdout bool
}

type Option func(config *Options)

func WithLevel(level string) Option {
	return func(o *Options) {
		o.level = level
	}
}

func WithPath(path string) Option {
	return func(o *Options) {
		o.path = path
	}
}

func WithMaxAge(days int) Option {
	return func(o *Options) {
		o.maxAge = days
	}
}

func WithMaxSize(megabytes int) Option {
	return func(o *Options) {
		o.maxSize = megabytes
	}
}

func WithMaxBackups(maxBackups int) Option {
	return func(o *Options) {
		o.maxBackups = maxBackups
	}
}

func WithCompress(compress bool) Option {
	return func(o *Options) {
		o.compress = compress
	}
}

func WithFormat(format string) Option {
	return func(o *Options) {
		o.format = format
	}
}

func WithEnableStdout(enableStdout bool) Option {
	return func(o *Options) {
		o.enableStdout = enableStdout
	}
}

// Init initializes the logger.
func Init(options ...Option) {
	opts := &Options{
		level:        "info",
		path:         "./logs/app.log",
		maxAge:       7,
		maxSize:      128,
		maxBackups:   32,
		compress:     true,
		format:       "text",
		enableStdout: true,
	}
	for _, option := range options {
		option(opts)
	}
	var writer io.Writer
	writer = &lumberjack.Logger{
		Filename:   opts.path,
		MaxSize:    opts.maxSize,
		MaxAge:     opts.maxAge,
		MaxBackups: opts.maxBackups,
		Compress:   opts.compress,
	}
	if opts.enableStdout {
		writer = io.MultiWriter(os.Stdout, writer)
	}
	handlerOptions := &slog.HandlerOptions{
		Level:     parseLevel(opts.level),
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				source.File = filepath.Base(source.File)
			}
			return a
		},
	}
	var handler slog.Handler
	switch strings.ToLower(opts.format) {
	case "json":
		handler = slog.NewJSONHandler(writer, handlerOptions)
	default:
		handler = slog.NewTextHandler(writer, handlerOptions)
	}
	slog.SetDefault(slog.New(&TraceContextHandler{
		Handler: handler,
	}))
	slog.Info("logger initialized", "path", opts.path, "format", opts.format, "level", opts.level)
}

// parseLevel parses a log level string and returns the corresponding slog.Level.
func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		slog.Warn("invalid log level, using info as default", "level", level)
		return slog.LevelInfo
	}
}
