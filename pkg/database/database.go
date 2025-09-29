/*
Copyright © 2025 lixw
*/
package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Options struct {
	url             string
	connMaxIdleTime time.Duration
	connMaxLifeTime time.Duration
	maxIdleConns    int
	maxOpenConns    int
	slowThreshold   time.Duration
	plugins         []gorm.Plugin
}

type Option func(*Options)

func WithUrl(url string) Option {
	return func(o *Options) {
		o.url = url
	}
}

func WithConnMaxIdleTime(connMaxIdleTime time.Duration) Option {
	return func(o *Options) {
		o.connMaxIdleTime = connMaxIdleTime
	}
}

func WithConnMaxLifeTime(connMaxLifeTime time.Duration) Option {
	return func(o *Options) {
		o.connMaxLifeTime = connMaxLifeTime
	}
}

func WithMaxIdleConns(maxIdleConns int) Option {
	return func(o *Options) {
		o.maxIdleConns = maxIdleConns
	}
}

func WithMaxOpenConns(maxOpenConns int) Option {
	return func(o *Options) {
		o.maxOpenConns = maxOpenConns
	}
}

func WithSlowThreshold(slowThreshold time.Duration) Option {
	return func(o *Options) {
		o.slowThreshold = slowThreshold
	}
}

func WithPlugins(plugins ...gorm.Plugin) Option {
	return func(o *Options) {
		o.plugins = append(o.plugins, plugins...)
	}
}

type gormLogger struct {
	slowThreshold time.Duration
	level         slog.Level
	logHandler    slog.Handler
}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	// 根据GORM的日志级别调整slog级别（简化实现）
	switch level {
	case logger.Silent:
		newLogger.level = slog.LevelError + 1
	case logger.Error:
		newLogger.level = slog.LevelError
	case logger.Warn:
		newLogger.level = slog.LevelWarn
	case logger.Info:
		newLogger.level = slog.LevelInfo
	}
	return &newLogger
}

func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.writeLog(ctx, slog.LevelInfo, msg, data...)
}

func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.writeLog(ctx, slog.LevelWarn, msg, data...)
}
func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	// 记录错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.writeLog(ctx, slog.LevelWarn, "%s [%d rows in %v]", sql, rows, elapsed)
		} else {
			l.writeLog(ctx, slog.LevelError, "%s [%d rows in %v] failed: %v", sql, rows, elapsed, err)
		}
		return
	}

	// 记录慢查询
	if l.slowThreshold != 0 && elapsed > l.slowThreshold {
		l.writeLog(ctx, slog.LevelWarn, "%s [%d rows in %v]", sql, rows, elapsed)
		return
	}

	// 常规日志
	l.writeLog(ctx, slog.LevelInfo, "%s [%d rows in %v]", sql, rows, elapsed)
}

func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.writeLog(ctx, slog.LevelError, msg, data...)
}

func (l *gormLogger) writeLog(ctx context.Context, level slog.Level, content string, args ...interface{}) {
	if level < l.level {
		return
	}
	var pc uintptr
	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(5, pcs[:])
	pc = pcs[0]
	r := slog.NewRecord(time.Now(), level, fmt.Sprintf(content, args...), pc)
	_ = l.logHandler.Handle(ctx, r)
}

func New(options ...Option) (*gorm.DB, error) {
	opts := &Options{
		maxIdleConns:    5,
		maxOpenConns:    10,
		connMaxIdleTime: 5 * time.Minute,
		connMaxLifeTime: 10 * time.Minute,
		slowThreshold:   time.Millisecond * 500,
	}
	for _, option := range options {
		option(opts)
	}
	if opts.url == "" {
		return nil, fmt.Errorf("database url is not set (use WithUrl to configure)")
	}
	db, err := gorm.Open(mysql.Open(opts.url),
		&gorm.Config{
			Logger: &gormLogger{
				slowThreshold: opts.slowThreshold,
				level:         slog.LevelInfo,
				logHandler:    slog.Default().Handler(),
			},
			DisableForeignKeyConstraintWhenMigrating: true,
		})
	if err != nil {
		return nil, err
	}
	for _, plugin := range opts.plugins {
		if err := db.Use(plugin); err != nil {
			return nil, err
		}
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxIdleTime(opts.connMaxIdleTime)
	sqlDB.SetConnMaxLifetime(opts.connMaxLifeTime)
	sqlDB.SetMaxIdleConns(opts.maxIdleConns)
	sqlDB.SetMaxOpenConns(opts.maxOpenConns)
	slog.Info("database create successfully")
	return db, nil
}

type DatabaseService struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *DatabaseService {
	return &DatabaseService{db: db}
}

func (s *DatabaseService) Start(ctx context.Context) error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	err = sqlDB.PingContext(ctx)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "database connect successfully")
	return nil
}

func (s *DatabaseService) Stop(ctx context.Context) error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "database disconnect successfully")
	return nil
}
