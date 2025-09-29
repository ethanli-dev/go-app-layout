/*
Copyright © 2025 lixw
*/
package i18n

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bytedance/sonic"
)

type localeKey struct{}

var (
	defaultLocale string
	locales       = make(map[string]map[string]string)
	mu            sync.RWMutex
)

type Options struct {
	lang     string
	fs       fs.FS
	basePath string
}

type Option func(*Options)

func WithLang(lang string) Option {
	return func(o *Options) { o.lang = lang }
}

func WithEmbedFS(efs fs.FS, basePath string) Option {
	return func(o *Options) {
		o.fs = efs
		o.basePath = basePath
	}
}

func WithStaticFS(path string) Option {
	return func(o *Options) {
		o.fs = os.DirFS(path)
		o.basePath = "."
	}
}

// Init initializes the i18n
func Init(options ...Option) error {
	opts := &Options{
		basePath: ".",
		fs:       os.DirFS("."),
		lang:     "zh-CN",
	}
	for _, option := range options {
		option(opts)
	}
	defaultLocale = opts.lang
	if defaultLocale == "" {
		return errors.New("default language cannot be empty")
	}
	err := fs.WalkDir(opts.fs, opts.basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			slog.Warn("walk dir error", "path", path, "err", err)
			return nil
		}
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), ".json") {
			return nil
		}
		buf, err := fs.ReadFile(opts.fs, path)
		if err != nil {
			return err
		}
		var dict map[string]string
		if err := sonic.Unmarshal(buf, &dict); err != nil {
			return err
		}
		lang := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		mu.Lock()
		locales[lang] = dict
		mu.Unlock()
		return nil
	})
	mu.RLock()
	_, ok := locales[defaultLocale]
	mu.RUnlock()
	if !ok {
		return fmt.Errorf("default language %s not found", defaultLocale)
	}

	slog.Info("i18n initialized", "path", opts.basePath)
	return err
}

func SetLocale(ctx context.Context, locale string) context.Context {
	if locale == "" {
		mu.RLock()
		locale = defaultLocale
		mu.RUnlock()
	}
	return context.WithValue(ctx, localeKey{}, locale)
}

func GetLocale(ctx context.Context) string {
	locale := ctx.Value(localeKey{})
	if s, ok := locale.(string); ok && s != "" {
		return s
	}
	mu.RLock()
	defer mu.RUnlock()
	return defaultLocale
}

func Translate(locale string, code any, args ...any) string {
	key := fmt.Sprintf("%v", code)
	mu.RLock()
	defer mu.RUnlock()
	// 1. 尝试从指定语言中获取
	if dict, ok := locales[locale]; ok {
		if msg, exist := dict[key]; exist {
			return fmt.Sprintf(msg, args...)
		}
	}
	// 2. 尝试从默认语言中获取
	if dict, ok := locales[defaultLocale]; ok {
		if msg, exist := dict[key]; exist {
			return fmt.Sprintf(msg, args...)
		}
	}
	// 3. 未找到翻译
	return fmt.Sprintf("unknown message [code=%s, locale=%s]", key, locale)
}

func Localize(ctx context.Context, code any, args ...any) string {
	return Translate(GetLocale(ctx), code, args)
}
