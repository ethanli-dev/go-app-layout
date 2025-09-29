//go:build wireinject
// +build wireinject

/*
Copyright © 2025 lixw
*/
package server

import (
	"github.com/ethanli-dev/go-app-layout/buildinfo"
	"github.com/ethanli-dev/go-app-layout/internal/router"
	"github.com/ethanli-dev/go-app-layout/locales"
	"github.com/ethanli-dev/go-app-layout/pkg/app"
	"github.com/ethanli-dev/go-app-layout/pkg/config"
	"github.com/ethanli-dev/go-app-layout/pkg/database"
	"github.com/ethanli-dev/go-app-layout/pkg/i18n"
	"github.com/ethanli-dev/go-app-layout/pkg/logging"
	"github.com/ethanli-dev/go-app-layout/pkg/web"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func createRoutes(cfg *config.Config, db *gorm.DB) (*router.Router, error) {
	panic(wire.Build(router.New))
}

func CreateApp(configPath string) (*app.App, error) {
	cfg, err := config.New(configPath)
	if err != nil {
		return nil, err
	}
	var logOpts []logging.Option
	if cfg.Logging != nil {
		logOpts = []logging.Option{
			logging.WithLevel(cfg.Logging.Level),
			logging.WithPath(cfg.Logging.Path),
			logging.WithMaxAge(cfg.Logging.MaxAge),
			logging.WithMaxSize(cfg.Logging.MaxSize),
			logging.WithMaxBackups(cfg.Logging.MaxBackups),
			logging.WithCompress(cfg.Logging.Compress),
			logging.WithFormat(cfg.Logging.Format),
			logging.WithEnableStdout(buildinfo.IsDev()),
		}
	}
	logging.Init(logOpts...)
	var i18nOpts []i18n.Option
	if cfg.Server != nil {
		i18nOpts = []i18n.Option{
			i18n.WithLang(cfg.Server.Locale),
			i18n.WithEmbedFS(locales.I18nFS, "."),
		}
	}
	if err := i18n.Init(i18nOpts...); err != nil {
		return nil, err
	}
	var dbOpts []database.Option
	if cfg.Database != nil {
		dbOpts = []database.Option{
			database.WithUrl(cfg.Database.Url),
			database.WithConnMaxIdleTime(cfg.Database.ConnMaxIdleTime),
			database.WithConnMaxLifeTime(cfg.Database.ConnMaxLifeTime),
			database.WithMaxIdleConns(cfg.Database.MaxIdleConns),
			database.WithMaxOpenConns(cfg.Database.MaxOpenConns),
			database.WithSlowThreshold(cfg.Database.SlowThreshold),
		}
	}
	db, err := database.New(dbOpts...)
	if err != nil {
		return nil, err
	}
	routes, err := createRoutes(cfg, db)
	if err != nil {
		return nil, err
	}
	var webOpts []web.Option
	if cfg.Server != nil {
		webOpts = []web.Option{
			web.WithAddress(cfg.Server.Addr),
			web.WithBasePath(cfg.Server.BasePath),
			web.WithReadTimeout(cfg.Server.ReadTimeout),
			web.WithWriteTimeout(cfg.Server.WriteTimeout),
			web.WithIdleTimeout(cfg.Server.IdleTimeout),
			web.WithMaxHeaderBytes(cfg.Server.MaxHeaderBytes),
		}
	}
	webServer := web.New(webOpts...).UseRoutes(routes.Register)
	var appOpts []app.Option
	if cfg.Server != nil {
		appOpts = []app.Option{
			app.WithStartTimeout(cfg.Server.StartTimeout),
			app.WithShutdownTimeout(cfg.Server.ShutdownTimeout),
		}
	}

	return app.New(appOpts...).Use(database.NewService(db), webServer), nil
}
