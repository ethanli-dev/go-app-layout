/*
Copyright © 2025 lixw
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/ethanli-dev/go-app-layout/buildinfo"
	"github.com/ethanli-dev/go-app-layout/cmd/server"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "go-app-layout",
	Version: buildinfo.Short(),
	Short:   "Go App Layout - 基于 Go 语言的应用程序开发框架",
	Long: `Go App Layout 是一个集成了配置管理、日志、数据库等组件的Go应用开发框架。
支持多环境配置、依赖注入和标准化的项目结构，适用于快速开发各类后端服务。

使用示例:
  # 启动服务
  # 启动服务（两种等效方式）
  go-app-layout server
  go-app-layout start

  # 查看版本信息
  go-app-layout --version`,
	SilenceUsage:  true,
	SilenceErrors: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		slog.Error("error executing command", "err", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(server.StartCmd)
}
