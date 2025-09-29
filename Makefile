.DEFAULT_GOAL := help

# 检查命令是否存在的函数
check-command = $(if $(shell command -v $(1) 2>/dev/null),,$(error 未找到命令: $(1)，请先安装))

# 声明伪目标，避免与同名文件冲突
.PHONY: help gen clean build run

help:
	@echo "Usage: make [target]"
	@echo "可用目标:"
	@echo "  gen          生成Wire依赖注入代码"
	@echo "  clean        清理构建产物和临时文件"
	@echo "  build        编译项目生成可执行文件（依赖clean）"
	@echo "  run          直接运行项目（开发模式，依赖gen）"

gen:
	@$(call check-command, bash)  # 检查bash是否存在
	@if [ ! -f "scripts/build.sh" ]; then \
		echo "错误：未找到脚本文件 scripts/build.sh"; \
		exit 1; \
	fi
	@bash scripts/build.sh gen

clean:
	@$(call check-command, bash)  # 检查bash是否存在
	@if [ ! -f "scripts/build.sh" ]; then \
		echo "错误：未找到脚本文件 scripts/build.sh"; \
		exit 1; \
	fi
	@bash scripts/build.sh clean

build:
	@$(call check-command, bash)  # 检查bash是否存在
	@if [ ! -f "scripts/build.sh" ]; then \
		echo "错误：未找到脚本文件 scripts/build.sh"; \
		exit 1; \
	fi
	@bash scripts/build.sh build

run: gen
	@$(call check-command, bash)  # 检查bash是否存在
	@if [ ! -f "scripts/build.sh" ]; then \
		echo "错误：未找到脚本文件 scripts/build.sh"; \
		exit 1; \
	fi
	@bash scripts/build.sh run