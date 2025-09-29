FROM golang:1.24-alpine AS builder
LABEL authors="lixw"

# 安装构建依赖（git用于版本信息，bash用于执行脚本，tzdata支持时区）
RUN apk add --no-cache git bash tzdata

WORKDIR /app

# 处理依赖
COPY go.mod go.sum ./
RUN go mod download && go mod verify  # 验证依赖完整性

# 复制源代码
COPY . .

# 构建前检查脚本存在性，提前暴露错误
RUN if [ ! -f "./scripts/build.sh" ]; then \
        echo "错误：构建脚本 ./scripts/build.sh 不存在"; \
        exit 1; \
    fi && \
    chmod +x ./scripts/build.sh && \
    ./scripts/build.sh build

# 运行阶段：轻量级基础镜像
FROM alpine:3.22
LABEL authors="lixw"

# 配置时区（上海）并安装必要运行时依赖
ENV TZ=Asia/Shanghai
RUN apk add --no-cache ca-certificates tzdata && \
    rm -rf /var/cache/apk/*  # 清理缓存

WORKDIR /app

# 仅复制运行必需的文件
COPY --from=builder /app/bin/go-app-layout ./app
COPY --from=builder /app/config/prod.yml ./config/prod.yml

# 确保二进制文件可执行
RUN chmod 755 /app/app

# 暴露应用端口
EXPOSE 8080

# 支持通过环境变量动态切换配置文件
CMD ["./app", "server", "-c", "config/prod.yml"]
