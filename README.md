# Nav

一个适合个人使用的导航管理系统，包含前台导航页、后台管理面板、SQLite 数据存储，以及基于 Docker Compose 的轻量部署方式。

项目的目标很直接：把常用站点、工具和平台入口统一收拢，既能在前台快速查找和打开，也能在后台方便地维护菜单、应用、用户和站点设置。

## 功能

### 前台导航

1. 展示所有已启用的菜单和应用
2. 支持一级菜单和二级菜单分组展示
3. 支持按名称、描述、分组进行搜索
4. 支持目录导航和返回顶部
5. 支持已登录用户从前台直接进入后台管理

### 后台管理

1. 仪表盘总览
2. 应用管理
3. 菜单管理
4. 用户管理
5. 站点设置

## 技术栈

### 前端

1. Vue 3
2. Vue Router
3. Vuex
4. Element Plus
5. Axios
6. Iconify

### 后端

1. Go
2. Gin
3. Gorm
4. SQLite
5. JWT

### 部署

1. Docker
2. Docker Compose
3. Nginx

## 项目结构

```text
nav/
├─ backend/                 # Go 后端
│  ├─ config/               # 配置加载
│  ├─ controllers/          # 控制器
│  ├─ database/             # 数据库初始化
│  ├─ middlewares/          # 中间件
│  ├─ models/               # 数据模型
│  ├─ routes/               # 路由
│  ├─ services/             # 服务层
│  ├─ utils/                # 工具函数
│  ├─ .env                  # 后端环境变量
│  ├─ Dockerfile            # 后端 Dockerfile
│  └─ main.go               # 后端入口
├─ frontend/                # Vue 前端
│  ├─ src/
│  ├─ .env                  # 前端环境变量
│  ├─ Dockerfile            # 前端 Dockerfile
│  └─ nginx.conf            # 前端 Nginx 配置
├─ data/                    # SQLite 数据目录
├─ docker-compose.yml       # Docker Compose 配置
├─ go.mod
├─ go.sum
└─ README.md
```

## 数据库说明

项目当前仅使用 SQLite。

后端环境变量文件：

[backend/.env](/C:/Users/dasdw/Desktop/nav/backend/.env)

当前数据库配置：

```env
DB_PATH=./data/nav.db
JWT_SECRET=your-secret
```

这里要区分两种路径：

1. 程序内部读取路径是 `./data/nav.db`
2. 宿主机实际文件路径是根目录下的 `data/nav.db`

也就是说，在当前 Docker Compose 配置下，真正需要备份的数据库文件是：

```text
data/nav.db
```

## 默认账号

首次启动时如果数据库为空，系统会自动创建默认管理员账号：

```text
用户名：admin
密码：admin
```

首次登录后建议立即修改密码。

## 本地开发运行

### 启动后端

在项目根目录执行：

```bash
go run ./backend
```

后端默认监听：

```text
http://localhost:8080
```

### 启动前端

在 `frontend` 目录执行：

```bash
npm install
npm run serve
```

前端默认访问：

```text
http://localhost:8081
```

具体端口以 Vue CLI 实际输出为准。

## Docker 部署

当前 Docker 目录结构已经整理为：

1. 后端 Dockerfile 在 [backend/Dockerfile](/C:/Users/dasdw/Desktop/nav/backend/Dockerfile)
2. 前端 Dockerfile 在 [frontend/Dockerfile](/C:/Users/dasdw/Desktop/nav/frontend/Dockerfile)
3. 前端 Nginx 配置在 [frontend/nginx.conf](/C:/Users/dasdw/Desktop/nav/frontend/nginx.conf)
4. 根目录统一通过 [docker-compose.yml](/C:/Users/dasdw/Desktop/nav/docker-compose.yml) 启动

### 一键启动

在项目根目录执行：

```bash
docker compose up -d --build
```

启动后：

1. 前端访问地址：`http://localhost`
2. 后端接口地址：`http://localhost:8080`

停止服务：

```bash
docker compose down
```

重新构建并启动：

```bash
docker compose up -d --build
```

## Compose 配置说明

当前 [docker-compose.yml](/C:/Users/dasdw/Desktop/nav/docker-compose.yml) 做了这些事情：

1. 后端使用根目录作为构建上下文，并指定 `backend/Dockerfile`
2. 前端使用 `frontend` 目录作为构建上下文，并指定 `frontend/Dockerfile`
3. 后端暴露 `8080`
4. 前端暴露 `80`
5. 本机 `./data` 挂载到容器内 `/app/backend/data`

所以容器内程序读取的是：

```text
./data/nav.db
```

而宿主机真正落盘的是：

```text
data/nav.db
```

## 单独构建镜像

### 构建后端镜像

在项目根目录执行：

```bash
docker build -t nav-backend:latest -f backend/Dockerfile .
```

### 构建前端镜像

在项目根目录执行：

```bash
docker build -t nav-frontend:latest ./frontend
```

## 单独运行容器

### 创建网络

```bash
docker network create nav-network
```

### 启动后端

```bash
docker run -d \
  --name nav-backend \
  --network nav-network \
  -p 8080:8080 \
  -v $(pwd)/data:/app/backend/data \
  nav-backend:latest
```

### 启动前端

```bash
docker run -d \
  --name nav-frontend \
  --network nav-network \
  -p 80:80 \
  nav-frontend:latest
```

如果你在 Windows PowerShell 下运行，挂载路径可以写成：

```powershell
docker run -d `
  --name nav-backend `
  --network nav-network `
  -p 8080:8080 `
  -v "${PWD}\data:/app/backend/data" `
  nav-backend:latest
```

## 容器常用命令

### 查看运行状态

```bash
docker ps
```

### 查看后端日志

```bash
docker logs -f nav-backend
```

### 查看前端日志

```bash
docker logs -f nav-frontend
```

### 进入后端容器

```bash
docker exec -it nav-backend sh
```

### 进入前端容器

```bash
docker exec -it nav-frontend sh
```

## 前后端访问关系

前端环境变量文件：

[frontend/.env](/C:/Users/dasdw/Desktop/nav/frontend/.env)

当前配置：

```env
VUE_APP_API_URL=/api
```

也就是说：

1. 浏览器访问前端页面时，接口统一走 `/api`
2. Nginx 再把 `/api` 代理到后端服务 `nav-backend:8080`

## 适用场景

当前这套方案更适合：

1. 个人导航站
2. 单机部署
3. 轻量后台管理
4. 小规模内部自用

如果以后要做多实例部署、共享数据库、复杂权限体系，再考虑把 SQLite 升级为 MySQL 或 PostgreSQL 会更稳。

## 建议

1. 首次部署完成后先登录后台检查菜单和应用数据
2. 及时修改默认管理员密码
3. 定期备份 `data/nav.db`
4. 如果只是个人使用，当前这套 Docker + SQLite 已经足够轻量稳定
