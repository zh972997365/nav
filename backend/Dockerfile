# ========== 阶段 1：构建 Go 应用 ==========
FROM golang:1.24-alpine AS builder

# 设置国内代理
ENV GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0

# 设置工作目录
WORKDIR /src

# 下载依赖文件
COPY go.mod go.sum* ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# 复制源码
COPY backend ./backend

# 构建应用
WORKDIR /src/backend
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-s -w" -trimpath -o /out/nav-backend .

# ========== 阶段 2：运行环境 ==========
FROM alpine:3.20

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app/backend

# 构建产物
COPY --from=builder /out/nav-backend ./main
COPY backend/.env ./

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./main"]