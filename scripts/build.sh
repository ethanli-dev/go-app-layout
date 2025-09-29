#!/bin/bash
set -euo pipefail  # 启用严格模式，遇到错误立即退出，避免隐蔽问题

APP_NAME="go-app-layout"
BIN_DIR="bin"
MAIN_PATH="main.go"
BUILD_PACKAGE="github.com/ethanli-dev/go-app-layout/buildinfo"

if ! git rev-parse --is-inside-work-tree &> /dev/null; then
    echo "错误：当前目录不是Git仓库，请初始化仓库或检查路径"
    exit 1
fi

if GIT_TAG_TEMP=$(git describe --tags --abbrev=0 2>/dev/null); then
    GIT_TAG="${GIT_TAG:-$GIT_TAG_TEMP}"
else
    GIT_TAG="${GIT_TAG:-v0.0.0}"
fi

GIT_COMMIT="${GIT_COMMIT:-$(git rev-parse --short HEAD)}"
BUILD_DATE="$(date +%Y-%m-%dT%H:%M:%SZ)"
GO_VERSION="$(go version | awk '{print $3}')"
BUILD_MODE="release"

check_command() {
    if ! command -v "$1" &> /dev/null; then
        echo "错误：未找到命令 '$1'，请先安装"
        exit 1
    fi
}

clean() {
    echo "清理构建产物..."
    rm -rf "$BIN_DIR"
    find . -name "*.log" -delete
    echo "清理完成"
}

build() {
    export GO111MODULE=on
    export GOPROXY="${GOPROXY:-https://goproxy.cn}"

    echo "安装依赖..."
    if ! go mod download && go mod verify; then
        echo "错误：依赖安装或验证失败"
        exit 1
    fi

    echo "开始编译项目（版本: $GIT_TAG, 提交: $GIT_COMMIT）..."
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

main() {
    check_command "go"
    check_command "git"
    clean
    build
}

main
