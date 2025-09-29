/*
Copyright © 2025 lixw
*/
package buildinfo

import "fmt"

// 构建信息变量（未导出，防止运行时篡改）
var (
	name      = "Go App Layout" // 项目名称
	version   = "v0.0.1"        // 版本号（如 v1.2.3）
	date      = "unknown"       // 构建时间（如 2025-09-28T12:34:56Z）
	goVersion = "unknown"       // 构建使用的Go版本
	mode      = "dev"           // 运行模式（dev/release）
	commit    = "unknown"       // Git提交哈希（短格式）
)

// 运行模式类型定义，限制只能为 dev 或 release
type RunMode string

const (
	ModeDev     RunMode = "dev"     // 开发模式
	ModeRelease RunMode = "release" // 生产模式
)

// 以下为 getter 方法，防止直接修改内部变量

// Name 返回项目名称
func Name() string {
	return name
}

// Version 返回版本号
func Version() string {
	return version
}

// Date 返回构建时间
func Date() string {
	return date
}

// GoVersion 返回构建使用的Go版本
func GoVersion() string {
	return goVersion
}

// Mode 返回运行模式（dev/release）
func Mode() RunMode {
	return RunMode(mode)
}

// Commit 返回Git提交哈希
func Commit() string {
	return commit
}

// IsRelease 判断是否为生产模式
func IsRelease() bool {
	return Mode() == ModeRelease
}

// IsDev 判断是否为开发模式
func IsDev() bool {
	return Mode() == ModeDev
}

func String() string {
	return fmt.Sprintf("%s %s (%s) built with %s [mode: %s, commit: %s]",
		Name(), Version(), Date(), GoVersion(), Mode(), Commit())
}

func Short() string {
	return fmt.Sprintf("%s (%s) built with %s [mode: %s, commit: %s]",
		Version(), Date(), GoVersion(), Mode(), Commit())
}
