#!/bin/bash
set -euo pipefail

APP_NAME="go-app-layout"
BIN_DIR="bin"
MAIN_PATH="main.go"
BUILD_PACKAGE="github.com/ethanli-dev/go-app-layout/buildinfo"
DEFAULT_CONFIG="config/dev.yml"
WIRE_GEN_PATH="./cmd/server"
BUILD_MODE="release"

# 初始化变量
if ! git rev-parse --is-inside-work-tree &> /dev/null; then
    echo "错误：当前目录不是Git仓库，请初始化仓库或检查路径"
    exit 1
fi

# 标签和提交信息处理
GIT_TAG="$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")"
GIT_COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")"
BUILD_DATE="$(date +%Y-%m-%dT%H:%M:%SZ)"
GO_VERSION="$(go version | awk '{print $3}')"

# 检查必要命令
check_command() {
    if ! command -v "$1" &> /dev/null; then
        echo "错误：未找到命令 '$1'，请先安装"
        exit 1
    fi
}

# 清理函数
clean() {
    echo "清理构建产物..."
    rm -rf "$BIN_DIR"
    find . -name "*.log" -delete
    echo "清理完成"
}

# 构建函数
build() {
    clean  # 构建前先清理

    check_command "go"
    check_command "git"

    export GO111MODULE=on
    export GOPROXY="${GOPROXY:-https://goproxy.cn}"

    echo "安装依赖..."
    if ! go mod download && go mod verify; then
        echo "错误：依赖安装或验证失败"
        exit 1
    fi

    echo "开始编译项目（版本: '$GIT_TAG', 提交: '$GIT_COMMIT'）..."
    mkdir -p "$BIN_DIR"

    if ! CGO_ENABLED=0 go build -trimpath \
        -ldflags "-s -w \
        -X \"$BUILD_PACKAGE.name=$APP_NAME\" \
        -X \"$BUILD_PACKAGE.version=$GIT_TAG\" \
        -X \"$BUILD_PACKAGE.date=$BUILD_DATE\" \
        -X \"$BUILD_PACKAGE.goVersion=$GO_VERSION\" \
        -X \"$BUILD_PACKAGE.mode=$BUILD_MODE\" \
        -X \"$BUILD_PACKAGE.commit=$GIT_COMMIT\" \
    " -o "$BIN_DIR/$APP_NAME" "$MAIN_PATH"; then
        echo "错误：项目构建失败"
        exit 1
    fi

    echo "编译完成，可执行文件位于: $BIN_DIR/$APP_NAME"
}

# 运行函数
# 支持传递额外参数（如配置文件路径），例如：./scripts/build.sh run -c config/prod.yml
run() {
    check_command "go"  # 确保go命令可用

    export GO111MODULE=on
    export GOPROXY="${GOPROXY:-https://goproxy.cn}"

    echo "安装依赖..."
    if ! go mod download && go mod verify; then
        echo "错误：依赖安装或验证失败"
        exit 1
    fi

    # 执行go run，传递所有额外参数（默认使用开发环境配置）
    echo "开始运行程序（开发模式）..."
    if [ $# -eq 0 ]; then
        # 无额外参数时，使用默认配置
        go run "$MAIN_PATH" server -c "$DEFAULT_CONFIG"
    else
        # 有额外参数时，优先使用用户传入的参数
        go run "$MAIN_PATH" "$@"
    fi
}

# gen函数（Wire代码生成）
gen() {
    # 确保GOPATH/bin在PATH中（兼容未显式设置GOPATH的情况）
    export GOPATH="${GOPATH:-$(go env GOPATH)}"
    export PATH="$PATH:$GOPATH/bin"

    # 检查wire是否已安装
    if ! command -v wire &> /dev/null; then
        echo "未找到wire工具，正在安装..."

        # 检查Go环境是否存在
        check_command "go"  # 复用现有检查函数

        # 安装wire
        echo "执行安装命令: go install github.com/google/wire/cmd/wire@latest"
        go install github.com/google/wire/cmd/wire@latest

        # 验证安装结果
        if ! command -v wire &> /dev/null; then
            echo "错误：wire安装失败，请检查GOPATH是否正确配置且在PATH中"
            echo "当前GOPATH: $GOPATH"
            echo "当前PATH: $PATH"
            exit 1
        fi
        echo "wire安装成功"
    fi

    # 生成Wire依赖注入代码
    echo "开始生成Wire代码: wire gen $WIRE_GEN_PATH"
    wire gen "$WIRE_GEN_PATH"

    echo "Wire代码生成完成"
}

# 命令参数解析（支持clean、build、run）
if [ $# -eq 0 ]; then
    # 无参数时默认执行build
    build
else
    cmd="$1"
    shift  # 移除命令名，保留后续参数（传给具体函数）
    case "$cmd" in
        clean)
            clean  # 执行清理（不接收额外参数）
            ;;
        build)
            build  # 执行构建（不接收额外参数）
            ;;
        run)
            run "$@"  # 传递所有额外参数给run函数
            ;;
        gen)
            gen  # 执行gen wire生成
            ;;
        *)
            echo "错误：仅支持 'clean'、'build'、'run' 或 'gen' 命令"
            exit 1
            ;;
    esac
fi
