# ========== 阶段 1：构建前端应用 ==========
FROM node:18-alpine as build

# 设置国内代理
RUN npm config set registry https://registry.npmmirror.com

# 设置工作目录
WORKDIR /app

# 下载依赖文件
COPY package*.json ./
RUN --mount=type=cache,target=/root/.npm \
    npm ci --no-audit --no-fund || \
    npm install --no-audit --no-fund

# 复制源码
COPY . .

# 构建应用
RUN npm run build

# ========== 阶段 2：运行环境 ==========
FROM nginx:alpine

# Nginx配置
COPY nginx.conf /etc/nginx/nginx.conf

# 构建产物
COPY --from=build /app/dist /usr/share/nginx/html

# 暴露端口
EXPOSE 80

# 启动 Nginx
CMD ["nginx", "-g", "daemon off;"]