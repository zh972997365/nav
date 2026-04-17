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

[backend/.env]

当前数据库配置：

```env
DB_PATH=./data/nav.db
JWT_SECRET=your-secret
```

备份的数据库文件是：

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

1. 后端 Dockerfile 在 [backend/Dockerfile]
2. 前端 Dockerfile 在 [frontend/Dockerfile]
3. 前端 Nginx 配置在 [frontend/nginx.conf]
4. 根目录统一通过 [docker-compose.yml] 启动

### 一键启动

在项目根目录执行：

```bash
docker compose up -d --build
```

启动后：

1. 前端访问地址：`http://localhost`
2. 后端接口地址：`http://localhost:8080`

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
