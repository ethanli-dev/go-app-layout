#!/bin/bash

# 确保GOPATH/bin在PATH中（兼容未显式设置GOPATH的情况）
export GOPATH="${GOPATH:-$(go env GOPATH)}"
export PATH="$PATH:$GOPATH/bin"

# 检查wire是否已安装
if ! command -v wire &> /dev/null; then
    echo "未找到wire工具，正在安装..."

    # 检查Go环境是否存在
    if ! command -v go &> /dev/null; then
        echo "错误：未找到go命令，请先安装Go环境（1.16+）"
        exit 1
    fi

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
echo "开始生成Wire代码: wire gen ./cmd/server"
wire gen ./cmd/server

echo "Wire代码生成完成"
